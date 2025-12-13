package oauth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/rs/zerolog"
)

// contextKey is used for storing values in context.
type contextKey string

const (
	// tokenInfoKey is the context key for TokenInfo.
	tokenInfoKey contextKey = "oauth:token_info"
)

// TokenInfoFromContext retrieves the TokenInfo from the request context.
func TokenInfoFromContext(ctx context.Context) (*TokenInfo, bool) {
	info, ok := ctx.Value(tokenInfoKey).(*TokenInfo)
	return info, ok
}

// ContextWithTokenInfo returns a new context with the TokenInfo stored.
func ContextWithTokenInfo(ctx context.Context, info *TokenInfo) context.Context {
	return context.WithValue(ctx, tokenInfoKey, info)
}

// BearerAuthMiddleware creates HTTP middleware that validates bearer tokens.
// It implements RFC 6750 Bearer Token Usage and returns proper WWW-Authenticate
// headers on failure per the MCP authorization specification.
type BearerAuthMiddleware struct {
	verifier    TokenVerifier
	config      *Config
	logger      *zerolog.Logger
	metadataURL string
}

// NewBearerAuthMiddleware creates a new bearer token authentication middleware.
func NewBearerAuthMiddleware(cfg *Config, verifier TokenVerifier, logger *zerolog.Logger) *BearerAuthMiddleware {
	return &BearerAuthMiddleware{
		verifier:    verifier,
		config:      cfg,
		logger:      logger,
		metadataURL: GetMetadataURL(cfg.ResourceServerURL),
	}
}

// Wrap wraps an HTTP handler with bearer token authentication.
func (m *BearerAuthMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip authentication for the metadata endpoint
		if r.URL.Path == MetadataPath {
			next.ServeHTTP(w, r)
			return
		}

		// Extract bearer token from Authorization header
		token, err := extractBearerToken(r)
		if err != nil {
			m.unauthorized(w, "invalid_request", err.Error())
			return
		}

		if token == "" {
			m.unauthorized(w, "invalid_request", "Bearer token required")
			return
		}

		// Verify the token
		info, err := m.verifier.Verify(r.Context(), token)
		if err != nil {
			if m.logger != nil {
				m.logger.Debug().Err(err).Msg("token verification failed")
			}
			m.unauthorized(w, "invalid_token", "Token verification failed")
			return
		}

		// Log successful authentication
		if m.logger != nil {
			m.logger.Debug().
				Str("subject", info.Subject).
				Str("client_id", info.ClientID).
				Strs("scopes", info.Scopes).
				Msg("authenticated request")
		}

		// Add token info to context
		ctx := ContextWithTokenInfo(r.Context(), info)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// unauthorized sends a 401 response with proper WWW-Authenticate header.
// This follows RFC 6750 and the MCP authorization specification.
func (m *BearerAuthMiddleware) unauthorized(w http.ResponseWriter, errorCode, description string) {
	// Build WWW-Authenticate header per RFC 6750 and MCP spec
	authHeader := fmt.Sprintf(
		`Bearer realm="mcp", resource_metadata="%s"`,
		m.metadataURL,
	)

	if errorCode != "" {
		authHeader += fmt.Sprintf(`, error="%s"`, errorCode)
	}
	if description != "" {
		authHeader += fmt.Sprintf(`, error_description="%s"`, description)
	}

	w.Header().Set("WWW-Authenticate", authHeader)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)

	// Return JSON error response
	_, _ = fmt.Fprintf(w, `{"error":"%s","error_description":"%s"}`, errorCode, description)
}

// extractBearerToken extracts the bearer token from the Authorization header.
func extractBearerToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", nil
	}

	// Must be "Bearer <token>"
	const prefix = "Bearer "
	if !strings.HasPrefix(authHeader, prefix) {
		return "", fmt.Errorf("authorization header must use Bearer scheme")
	}

	token := strings.TrimPrefix(authHeader, prefix)
	token = strings.TrimSpace(token)

	if token == "" {
		return "", fmt.Errorf("bearer token is empty")
	}

	return token, nil
}

// RequireScopes returns middleware that requires specific scopes.
func RequireScopes(scopes ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			info, ok := TokenInfoFromContext(r.Context())
			if !ok {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}

			if !info.HasAllScopes(scopes) {
				http.Error(w, "Insufficient scope", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
