// Package auth implements GitHub OAuth-backed authentication for the dashboard.
//
// Session model: a signed cookie carries a small payload (subject, issued_at,
// expires_at). Signing uses HMAC-SHA256 with a server-side secret. There is no
// server-side session store — the cookie itself is the session.
//
// Cookie format: base64url(payload-json) + "." + base64url(hmac-sha256(payload, secret))
package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/goccy/go-json"
)

// CookieName is the HTTP cookie that carries the dashboard session.
const CookieName = "dmarcguard_session"

// StateCookieName is the short-lived cookie that carries the OAuth state token
// used to defeat CSRF on the OAuth callback. It is deleted as soon as the
// callback validates it.
const StateCookieName = "dmarcguard_oauth_state"

// DefaultSessionTTL is used when AuthConfig.SessionTTLDays is unset (0).
const DefaultSessionTTL = 7 * 24 * time.Hour

// Errors returned by session decode. Callers treat any of these as
// "user is not authenticated" rather than distinguishing failure modes.
var (
	ErrSessionMissing   = errors.New("session cookie missing")
	ErrSessionMalformed = errors.New("session cookie malformed")
	ErrSessionBadSig    = errors.New("session cookie signature mismatch")
	ErrSessionExpired   = errors.New("session expired")
)

// Session is what we store inside the cookie. It is intentionally tiny — just
// enough to identify the user and enforce expiry. The allowlist check happens
// at session-creation time (callback handler), so by the time a session
// exists, the user is already authorized.
type Session struct {
	Subject   string `json:"sub"`   // GitHub username
	Email     string `json:"email"` // verified primary email
	IssuedAt  int64  `json:"iat"`
	ExpiresAt int64  `json:"exp"`
}

// SessionSigner signs and verifies session cookies using HMAC-SHA256.
// Construct via NewSessionSigner; key is the raw secret bytes (≥32 bytes).
type SessionSigner struct {
	key []byte
	ttl time.Duration
}

// NewSessionSigner constructs a signer from a base64-encoded secret string and
// a TTL. ttlDays<=0 selects DefaultSessionTTL.
func NewSessionSigner(b64Secret string, ttlDays int) (*SessionSigner, error) {
	key, err := base64.StdEncoding.DecodeString(b64Secret)
	if err != nil {
		// Fall back to raw bytes — let users paste a non-base64 secret if they want.
		key = []byte(b64Secret)
	}
	if len(key) < 32 {
		return nil, fmt.Errorf("session secret must be ≥32 bytes (got %d)", len(key))
	}
	ttl := DefaultSessionTTL
	if ttlDays > 0 {
		ttl = time.Duration(ttlDays) * 24 * time.Hour
	}
	return &SessionSigner{key: key, ttl: ttl}, nil
}

// TTL returns the configured session TTL.
func (s *SessionSigner) TTL() time.Duration { return s.ttl }

// Sign builds a session cookie value for the given identity.
func (s *SessionSigner) Sign(subject, email string, now time.Time) (string, error) {
	sess := Session{
		Subject:   subject,
		Email:     email,
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(s.ttl).Unix(),
	}
	payload, err := json.Marshal(sess)
	if err != nil {
		return "", err
	}
	encPayload := base64.RawURLEncoding.EncodeToString(payload)
	mac := hmacSum(s.key, []byte(encPayload))
	encSig := base64.RawURLEncoding.EncodeToString(mac)
	return encPayload + "." + encSig, nil
}

// Verify decodes a cookie value, checks the signature, and returns the
// session. Returns ErrSession* sentinels on every failure path.
func (s *SessionSigner) Verify(cookieValue string, now time.Time) (*Session, error) {
	if cookieValue == "" {
		return nil, ErrSessionMissing
	}
	dot := strings.IndexByte(cookieValue, '.')
	if dot < 0 || dot == len(cookieValue)-1 {
		return nil, ErrSessionMalformed
	}
	encPayload := cookieValue[:dot]
	encSig := cookieValue[dot+1:]

	gotSig, err := base64.RawURLEncoding.DecodeString(encSig)
	if err != nil {
		return nil, ErrSessionMalformed
	}
	expectedSig := hmacSum(s.key, []byte(encPayload))
	if subtle.ConstantTimeCompare(gotSig, expectedSig) != 1 {
		return nil, ErrSessionBadSig
	}

	payload, err := base64.RawURLEncoding.DecodeString(encPayload)
	if err != nil {
		return nil, ErrSessionMalformed
	}
	var sess Session
	if err := json.Unmarshal(payload, &sess); err != nil {
		return nil, ErrSessionMalformed
	}
	if now.Unix() >= sess.ExpiresAt {
		return nil, ErrSessionExpired
	}
	return &sess, nil
}

// SignState produces a short-lived signed token for the OAuth state parameter.
// Format: base64(random16).base64(hmac).base64(expiresUnix). 10-min lifetime.
func (s *SessionSigner) SignState(now time.Time) (string, error) {
	nonce := make([]byte, 16)
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	exp := now.Add(10 * time.Minute).Unix()
	body := fmt.Sprintf("%s.%d", base64.RawURLEncoding.EncodeToString(nonce), exp)
	mac := hmacSum(s.key, []byte(body))
	return body + "." + base64.RawURLEncoding.EncodeToString(mac), nil
}

// VerifyState checks a state token's signature and expiry.
func (s *SessionSigner) VerifyState(token string, now time.Time) error {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return ErrSessionMalformed
	}
	body := parts[0] + "." + parts[1]
	gotSig, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return ErrSessionMalformed
	}
	expectedSig := hmacSum(s.key, []byte(body))
	if subtle.ConstantTimeCompare(gotSig, expectedSig) != 1 {
		return ErrSessionBadSig
	}
	var exp int64
	if _, err := fmt.Sscanf(parts[1], "%d", &exp); err != nil {
		return ErrSessionMalformed
	}
	if now.Unix() >= exp {
		return ErrSessionExpired
	}
	return nil
}

func hmacSum(key, msg []byte) []byte {
	m := hmac.New(sha256.New, key)
	m.Write(msg)
	return m.Sum(nil)
}

// GenerateSessionSecret returns a base64-encoded 32-byte random string suitable
// for the AuthConfig.SessionSecret field. Used by --gen-session-secret.
func GenerateSessionSecret() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}
