package oauth

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/oauth2"
)

// IMAPScope is the only Google scope that grants XOAUTH2 access to IMAP.
// Narrower Gmail API scopes (gmail.readonly, gmail.modify) do not work over
// IMAP — Gmail rejects them with AUTHENTICATE errors. Microsoft 365's
// equivalent will be IMAP.AccessAsUser.All when that provider is added.
const IMAPScope = "https://mail.google.com/"

// Provider abstracts an OAuth2 identity provider that supports the device
// authorization grant and produces tokens usable as XOAUTH2 SASL credentials.
type Provider interface {
	// Endpoint returns OAuth2 authorization/token endpoints for use with
	// golang.org/x/oauth2.
	Endpoint() oauth2.Endpoint

	// DeviceAuthURL returns the device authorization endpoint
	// (RFC 8628 §3.1). x/oauth2 v0.36 has DeviceAuth helpers but its Config
	// reads this from a separate field; we expose it here for clarity.
	DeviceAuthURL() string

	// UserinfoURL returns the OpenID Connect userinfo endpoint, used by
	// --oauth-login to confirm which mailbox identity the token belongs to.
	UserinfoURL() string

	// Scopes returns the OAuth scopes to request.
	Scopes() []string
}

// ProviderByName returns the Provider for a given identifier from the config.
func ProviderByName(name string) (Provider, error) {
	switch name {
	case "google":
		return Google{}, nil
	default:
		return nil, fmt.Errorf("unknown oauth provider %q (supported: google)", name)
	}
}

// Google implements Provider for Google identity.
type Google struct{}

// Endpoint returns Google's OAuth2 endpoints.
func (Google) Endpoint() oauth2.Endpoint {
	return oauth2.Endpoint{
		AuthURL:       "https://accounts.google.com/o/oauth2/v2/auth",
		DeviceAuthURL: "https://oauth2.googleapis.com/device/code",
		TokenURL:      "https://oauth2.googleapis.com/token",
		AuthStyle:     oauth2.AuthStyleInParams,
	}
}

// DeviceAuthURL returns Google's device authorization endpoint.
func (Google) DeviceAuthURL() string {
	return "https://oauth2.googleapis.com/device/code"
}

// UserinfoURL returns Google's OIDC userinfo endpoint.
func (Google) UserinfoURL() string {
	return "https://openidconnect.googleapis.com/v1/userinfo"
}

// Scopes returns the IMAP scope plus OIDC scopes required to read userinfo
// during the bootstrap flow.
func (Google) Scopes() []string {
	return []string{IMAPScope, "openid", "email"}
}

// IsTerminalAuthError reports whether err indicates the refresh token is no
// longer usable and human action is required (re-running --oauth-login).
// Distinguishes permanent failures (revoked grant, deleted account, scope
// changes) from transient ones (5xx, timeouts, network blips) that should
// resolve on their own.
func IsTerminalAuthError(err error) bool {
	if err == nil {
		return false
	}
	var rerr *oauth2.RetrieveError
	if errors.As(err, &rerr) {
		// RFC 6749 §5.2 codes that mean "this grant will never work again".
		switch rerr.ErrorCode {
		case "invalid_grant", "invalid_client", "unauthorized_client", "unsupported_grant_type":
			return true
		}
		// 4xx without a known code is treated as terminal too — provider is
		// telling us the request itself is wrong, retrying won't fix it.
		if rerr.Response != nil && rerr.Response.StatusCode >= 400 && rerr.Response.StatusCode < 500 {
			return true
		}
	}
	// Last-resort string match for providers that don't populate ErrorCode
	// cleanly. invalid_grant in particular is the canonical revocation signal.
	if strings.Contains(err.Error(), "invalid_grant") {
		return true
	}
	return false
}

// NewTokenSource builds a self-refreshing TokenSource backed by the given
// refresh token. The source caches the access token in process memory and
// refreshes it ~10s before expiry (default behavior of x/oauth2.ReuseTokenSource).
func NewTokenSource(ctx context.Context, p Provider, clientID, clientSecret, refreshToken string) oauth2.TokenSource {
	cfg := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     p.Endpoint(),
		Scopes:       p.Scopes(),
	}
	// An expired stub token with a refresh token forces the first call to hit
	// the refresh endpoint — we don't have an access token at startup.
	stub := &oauth2.Token{RefreshToken: refreshToken}
	return cfg.TokenSource(ctx, stub)
}
