package oauth

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
)

// TokenVerifier validates access tokens and extracts their information.
type TokenVerifier interface {
	// Verify validates the given token and returns its information.
	// Returns an error if the token is invalid or expired.
	Verify(ctx context.Context, token string) (*TokenInfo, error)
}

// OIDCVerifier validates tokens using OIDC/JWT verification.
type OIDCVerifier struct {
	provider  *oidc.Provider
	verifier  *oidc.IDTokenVerifier
	config    *Config
	initOnce  sync.Once
	initError error
}

// NewOIDCVerifier creates a new OIDC token verifier.
func NewOIDCVerifier(cfg *Config) *OIDCVerifier {
	return &OIDCVerifier{
		config: cfg,
	}
}

func (v *OIDCVerifier) init(ctx context.Context) error {
	v.initOnce.Do(func() {
		// Always use a custom HTTP client with a timeout.
		httpClient := &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: v.config.InsecureSkipVerify,
				},
			},
		}

		ctx = oidc.ClientContext(ctx, httpClient)

		provider, err := oidc.NewProvider(ctx, v.config.Issuer)
		if err != nil {
			v.initError = fmt.Errorf("failed to create OIDC provider: %w", err)
			return
		}
		v.provider = provider

		verifierConfig := &oidc.Config{
			ClientID:          v.config.Audience,
			SkipClientIDCheck: v.config.Audience == "",
			SkipIssuerCheck:   v.config.SkipIssuerCheck,
		}

		v.verifier = provider.Verifier(verifierConfig)
	})

	return v.initError
}

// Verify validates the given token using OIDC verification.
func (v *OIDCVerifier) Verify(ctx context.Context, token string) (*TokenInfo, error) {
	if err := v.init(ctx); err != nil {
		return nil, err
	}

	idToken, err := v.verifier.Verify(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("token verification failed: %w", err)
	}

	var claims struct {
		Subject   string      `json:"sub"`
		Audience  interface{} `json:"aud"`
		ClientID  string      `json:"azp"` // Authorized party (Keycloak/OAuth2)
		Scope     string      `json:"scope"`
		IssuedAt  int64       `json:"iat"`
		ExpiresAt int64       `json:"exp"`
	}

	if err := idToken.Claims(&claims); err != nil {
		return nil, fmt.Errorf("failed to parse token claims: %w", err)
	}

	// Parse audience (can be string or array)
	var audience []string
	switch aud := claims.Audience.(type) {
	case string:
		audience = []string{aud}
	case []interface{}:
		for _, a := range aud {
			if s, ok := a.(string); ok {
				audience = append(audience, s)
			}
		}
	default:
		return nil, fmt.Errorf("unexpected type for audience claim: %T", claims.Audience)
	}

	// Parse scopes (space-separated string)
	var scopes []string
	if claims.Scope != "" {
		scopes = strings.Split(claims.Scope, " ")
	}

	// Get all claims as extra data
	var extra map[string]interface{}
	_ = idToken.Claims(&extra)

	info := &TokenInfo{
		Subject:   claims.Subject,
		ClientID:  claims.ClientID,
		Scopes:    scopes,
		Audience:  audience,
		ExpiresAt: claims.ExpiresAt,
		IssuedAt:  claims.IssuedAt,
		Issuer:    idToken.Issuer,
		Extra:     extra,
	}

	// Validate audience if configured
	if v.config.Audience != "" && !info.HasAudience(v.config.Audience) {
		return nil, fmt.Errorf("token audience does not match expected audience: got %v, want %s",
			audience, v.config.Audience)
	}

	// Validate required scopes
	if len(v.config.RequiredScopes) > 0 && !info.HasAllScopes(v.config.RequiredScopes) {
		return nil, fmt.Errorf("token missing required scopes: has %v, requires %v",
			scopes, v.config.RequiredScopes)
	}

	return info, nil
}

// IntrospectionVerifier validates tokens using OAuth 2.0 Token Introspection (RFC 7662).
type IntrospectionVerifier struct {
	config     *Config
	httpClient *http.Client
}

// NewIntrospectionVerifier creates a new token introspection verifier.
func NewIntrospectionVerifier(cfg *Config) *IntrospectionVerifier {
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	if cfg.InsecureSkipVerify {
		httpClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	return &IntrospectionVerifier{
		config:     cfg,
		httpClient: httpClient,
	}
}

// introspectionResponse represents the response from the introspection endpoint.
type introspectionResponse struct {
	Active    bool        `json:"active"`
	Scope     string      `json:"scope"`
	ClientID  string      `json:"client_id"`
	Username  string      `json:"username"`
	TokenType string      `json:"token_type"`
	ExpiresAt int64       `json:"exp"`
	IssuedAt  int64       `json:"iat"`
	NotBefore int64       `json:"nbf"`
	Subject   string      `json:"sub"`
	Audience  interface{} `json:"aud"`
	Issuer    string      `json:"iss"`
	JTI       string      `json:"jti"`
}

// Verify validates the given token using token introspection.
func (v *IntrospectionVerifier) Verify(ctx context.Context, token string) (*TokenInfo, error) {
	if v.config.IntrospectionEndpoint == "" {
		return nil, errors.New("introspection endpoint not configured")
	}

	// Build request body
	data := url.Values{
		"token":         {token},
		"client_id":     {v.config.ClientID},
		"client_secret": {v.config.ClientSecret},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, v.config.IntrospectionEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create introspection request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := v.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("introspection request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("introspection returned status %d", resp.StatusCode)
	}

	var ir introspectionResponse
	if err := json.NewDecoder(resp.Body).Decode(&ir); err != nil {
		return nil, fmt.Errorf("failed to parse introspection response: %w", err)
	}

	if !ir.Active {
		return nil, errors.New("token is not active")
	}

	// Parse audience (can be string or array)
	var audience []string
	switch aud := ir.Audience.(type) {
	case string:
		audience = []string{aud}
	case []interface{}:
		for _, a := range aud {
			if s, ok := a.(string); ok {
				audience = append(audience, s)
			}
		}
	default:
		return nil, fmt.Errorf("unexpected type for audience claim: %T", ir.Audience)
	}

	// Parse scopes
	var scopes []string
	if ir.Scope != "" {
		scopes = strings.Split(ir.Scope, " ")
	}

	info := &TokenInfo{
		Subject:   ir.Subject,
		ClientID:  ir.ClientID,
		Scopes:    scopes,
		Audience:  audience,
		ExpiresAt: ir.ExpiresAt,
		IssuedAt:  ir.IssuedAt,
		Issuer:    ir.Issuer,
	}

	// Validate audience if configured
	if v.config.Audience != "" && !info.HasAudience(v.config.Audience) {
		return nil, fmt.Errorf("token audience does not match expected audience: got %v, want %s",
			audience, v.config.Audience)
	}

	// Validate required scopes
	if len(v.config.RequiredScopes) > 0 && !info.HasAllScopes(v.config.RequiredScopes) {
		return nil, fmt.Errorf("token missing required scopes: has %v, requires %v",
			scopes, v.config.RequiredScopes)
	}

	return info, nil
}

// CachingVerifier wraps a TokenVerifier with a cache to reduce validation overhead.
type CachingVerifier struct {
	verifier TokenVerifier
	cache    sync.Map
	ttl      time.Duration
}

type cacheEntry struct {
	info      *TokenInfo
	expiresAt time.Time
}

// NewCachingVerifier creates a verifier that caches validation results.
func NewCachingVerifier(verifier TokenVerifier, ttl time.Duration) *CachingVerifier {
	if ttl <= 0 {
		ttl = 5 * time.Minute // Default cache TTL
	}
	return &CachingVerifier{
		verifier: verifier,
		ttl:      ttl,
	}
}

// Verify validates the token, using cached results when available.
func (v *CachingVerifier) Verify(ctx context.Context, token string) (*TokenInfo, error) {
	// Check cache first
	if entry, ok := v.cache.Load(token); ok {
		if ce, ok := entry.(*cacheEntry); ok {
			if time.Now().Before(ce.expiresAt) {
				return ce.info, nil
			}
			// Cache entry expired, remove it
			v.cache.Delete(token)
		}
	}

	// Verify token
	info, err := v.verifier.Verify(ctx, token)
	if err != nil {
		return nil, err
	}

	// Cache the result
	// Use the lesser of token expiry or TTL
	expiresAt := time.Now().Add(v.ttl)
	if info.ExpiresAt > 0 {
		tokenExpiry := time.Unix(info.ExpiresAt, 0)
		if tokenExpiry.Before(expiresAt) {
			expiresAt = tokenExpiry
		}
	}

	v.cache.Store(token, &cacheEntry{
		info:      info,
		expiresAt: expiresAt,
	})

	return info, nil
}

// NewVerifier creates the appropriate token verifier based on configuration.
// If IntrospectionEndpoint is set, it uses introspection; otherwise, it uses OIDC.
func NewVerifier(cfg *Config) TokenVerifier {
	var verifier TokenVerifier

	if cfg.IntrospectionEndpoint != "" {
		verifier = NewIntrospectionVerifier(cfg)
	} else {
		verifier = NewOIDCVerifier(cfg)
	}

	// Wrap with caching
	return NewCachingVerifier(verifier, 5*time.Minute)
}
