// Package oauth provides OAuth 2.1 authentication for MCP servers.
// It implements the MCP authorization specification with support for:
// - OIDC/JWT token verification
// - Token introspection (RFC 7662)
// - Protected Resource Metadata (RFC 9728)
// - Bearer token middleware (RFC 6750)
package oauth

import (
	"errors"
	"net/url"
	"strings"
)

// Config holds OAuth2 configuration for the MCP server.
type Config struct {
	// Enabled determines whether OAuth2 authentication is active.
	Enabled bool

	// Issuer is the OAuth2/OIDC issuer URL (e.g., https://auth.example.com/realms/master).
	// Used for OIDC discovery and token validation.
	Issuer string

	// Audience is the expected audience claim in tokens. Usually the MCP server URL.
	// Required for token validation.
	Audience string

	// ClientID is the OAuth2 client ID for this MCP server.
	// Used for token introspection and optional client authentication.
	ClientID string

	// ClientSecret is the OAuth2 client secret for this MCP server.
	// Used for token introspection endpoint authentication.
	ClientSecret string

	// RequiredScopes are the scopes that must be present in the token.
	// If empty, no scope validation is performed.
	RequiredScopes []string

	// IntrospectionEndpoint is the URL for token introspection (RFC 7662).
	// If set, tokens will be validated via introspection instead of local JWT validation.
	IntrospectionEndpoint string

	// ResourceServerURL is the URL of this MCP server.
	// Used in Protected Resource Metadata for resource indicator validation.
	ResourceServerURL string

	// ResourceName is a human-readable name for this MCP server.
	ResourceName string

	// ResourceDocumentation is a URL to developer documentation.
	ResourceDocumentation string

	// SkipIssuerCheck disables issuer validation (for development only).
	SkipIssuerCheck bool

	// InsecureSkipVerify disables TLS certificate verification (for development only).
	InsecureSkipVerify bool
}

// Validate checks the configuration for required fields and consistency.
func (c *Config) Validate() error {
	if !c.Enabled {
		return nil
	}

	var errs []string

	if c.Issuer == "" {
		errs = append(errs, "issuer is required when OAuth is enabled")
	} else if _, err := url.Parse(c.Issuer); err != nil {
		errs = append(errs, "issuer must be a valid URL")
	}

	if c.Audience == "" {
		errs = append(errs, "audience is required when OAuth is enabled")
	}

	if c.ResourceServerURL == "" {
		errs = append(errs, "resource_server_url is required when OAuth is enabled")
	} else if _, err := url.Parse(c.ResourceServerURL); err != nil {
		errs = append(errs, "resource_server_url must be a valid URL")
	}

	if c.IntrospectionEndpoint != "" {
		if c.ClientID == "" || c.ClientSecret == "" {
			errs = append(errs, "client_id and client_secret are required when using introspection")
		}
		if _, err := url.Parse(c.IntrospectionEndpoint); err != nil {
			errs = append(errs, "introspection_endpoint must be a valid URL")
		}
	}

	if len(errs) > 0 {
		return errors.New("oauth config validation failed: " + strings.Join(errs, "; "))
	}

	return nil
}

// TokenInfo holds information extracted from a validated access token.
type TokenInfo struct {
	// Subject is the subject claim (typically user ID).
	Subject string

	// ClientID is the client that obtained the token.
	ClientID string

	// Scopes are the scopes granted to this token.
	Scopes []string

	// Audience is the intended audience for this token.
	Audience []string

	// ExpiresAt is the Unix timestamp when the token expires.
	ExpiresAt int64

	// IssuedAt is the Unix timestamp when the token was issued.
	IssuedAt int64

	// Issuer is the token issuer.
	Issuer string

	// Extra holds any additional claims from the token.
	Extra map[string]interface{}
}

// HasScope checks if the token has a specific scope.
func (t *TokenInfo) HasScope(scope string) bool {
	for _, s := range t.Scopes {
		if s == scope {
			return true
		}
	}
	return false
}

// HasAllScopes checks if the token has all the specified scopes.
func (t *TokenInfo) HasAllScopes(scopes []string) bool {
	for _, required := range scopes {
		if !t.HasScope(required) {
			return false
		}
	}
	return true
}

// HasAudience checks if the token is intended for a specific audience.
func (t *TokenInfo) HasAudience(audience string) bool {
	// Normalize both for comparison (remove trailing slashes)
	normalizedAudience := strings.TrimSuffix(audience, "/")
	for _, aud := range t.Audience {
		if strings.TrimSuffix(aud, "/") == normalizedAudience {
			return true
		}
	}
	return false
}
