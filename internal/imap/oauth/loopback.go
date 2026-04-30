package oauth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
	"strconv"
	"time"

	"golang.org/x/oauth2"
)

// LoopbackLogin runs the OAuth2 authorization code flow with PKCE against a
// localhost redirect URI (RFC 8252 §7.3 — the recommended flow for native /
// desktop apps). Returns a refresh token and the authenticated email.
//
// The flow:
//  1. Listen on 127.0.0.1 on a random free port.
//  2. Build the auth URL with redirect_uri=http://127.0.0.1:<port>/callback,
//     a CSRF state token, and a PKCE S256 code challenge.
//  3. Print the URL (and try to auto-open the user's browser).
//  4. Block until the callback handler receives ?code=... or ctx is cancelled.
//  5. Exchange the code (with the PKCE verifier) for tokens.
//  6. Fetch the userinfo email so the user knows which mailbox they authorized.
//
// Requires an OAuth client of type "Desktop app" — Google rejects loopback
// redirects on Web and TV/Limited Input client types.
func LoopbackLogin(ctx context.Context, p Provider, clientID, clientSecret string, prompt PromptFunc) (*DeviceLoginResult, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, fmt.Errorf("listen on loopback: %w", err)
	}
	defer func() { _ = listener.Close() }()

	port := listener.Addr().(*net.TCPAddr).Port
	redirectURI := "http://127.0.0.1:" + strconv.Itoa(port) + "/callback"

	state, err := randomToken(32)
	if err != nil {
		return nil, err
	}
	verifier, challenge, err := generatePKCE()
	if err != nil {
		return nil, err
	}

	cfg := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     p.Endpoint(),
		Scopes:       p.Scopes(),
		RedirectURL:  redirectURI,
	}
	authURL := cfg.AuthCodeURL(state,
		oauth2.AccessTypeOffline, // request refresh_token
		oauth2.ApprovalForce,     // force consent so a refresh_token is always returned
		oauth2.SetAuthURLParam("code_challenge", challenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
	)

	prompt(authURL, "")
	tryOpenBrowser(authURL)

	type callbackResult struct {
		code string
		err  error
	}
	resultCh := make(chan callbackResult, 1)

	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if e := q.Get("error"); e != "" {
			httpResponse(w, "Authorization failed: "+e+". You can close this tab.")
			resultCh <- callbackResult{err: fmt.Errorf("authorization error from provider: %s", e)}
			return
		}
		if got := q.Get("state"); got != state {
			httpResponse(w, "State mismatch — possible CSRF. You can close this tab.")
			resultCh <- callbackResult{err: errors.New("oauth state mismatch (possible CSRF or stale callback)")}
			return
		}
		code := q.Get("code")
		if code == "" {
			httpResponse(w, "Missing authorization code. You can close this tab.")
			resultCh <- callbackResult{err: errors.New("oauth callback missing code")}
			return
		}
		httpResponse(w, "Authorization complete. You can close this tab and return to the terminal.")
		resultCh <- callbackResult{code: code}
	})

	server := &http.Server{Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	go func() { _ = server.Serve(listener) }()
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdownCtx)
	}()

	var res callbackResult
	select {
	case res = <-resultCh:
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	if res.err != nil {
		return nil, res.err
	}

	tok, err := cfg.Exchange(ctx, res.code, oauth2.SetAuthURLParam("code_verifier", verifier))
	if err != nil {
		return nil, fmt.Errorf("exchange code for token: %w", err)
	}
	if tok.RefreshToken == "" {
		return nil, errors.New("provider returned no refresh token (tip: revoke prior consent at https://myaccount.google.com/permissions and retry)")
	}

	email, err := fetchUserinfoEmail(ctx, p, tok.AccessToken)
	if err != nil {
		email = ""
	}

	return &DeviceLoginResult{RefreshToken: tok.RefreshToken, Email: email}, nil
}

func httpResponse(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte("<!doctype html><html><body style=\"font-family:system-ui;padding:2rem\"><p>" + msg + "</p></body></html>"))
}

func randomToken(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// generatePKCE returns a verifier (43–128 chars) and its S256 challenge per
// RFC 7636. We use 32 random bytes which encode to a 43-char URL-safe string.
func generatePKCE() (verifier, challenge string, err error) {
	verifier, err = randomToken(32)
	if err != nil {
		return "", "", err
	}
	sum := sha256.Sum256([]byte(verifier))
	challenge = base64.RawURLEncoding.EncodeToString(sum[:])
	return verifier, challenge, nil
}

// tryOpenBrowser fires the OS browser at u. Best-effort — silently no-ops on
// failure, since the user can always copy the URL from the prompt.
func tryOpenBrowser(u string) {
	if _, err := url.Parse(u); err != nil {
		return
	}
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", u)
	case "linux":
		cmd = exec.Command("xdg-open", u)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", u)
	default:
		return
	}
	_ = cmd.Start()
}
