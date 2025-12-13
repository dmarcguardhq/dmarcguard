package oauth

import (
	"encoding/json"
	"net/http"
	"strings"
)

// ProtectedResourceMetadata represents OAuth 2.0 Protected Resource Metadata per RFC 9728.
// This metadata helps MCP clients discover how to authenticate with this server.
type ProtectedResourceMetadata struct {
	// Resource is the URI identifying the protected resource.
	Resource string `json:"resource"`

	// AuthorizationServers lists the authorization server issuer identifiers.
	AuthorizationServers []string `json:"authorization_servers,omitempty"`

	// JWKSURI is the URL to this resource server's JWK Set document.
	JWKSURI string `json:"jwks_uri,omitempty"`

	// ScopesSupported lists the scopes that may be requested for this resource.
	ScopesSupported []string `json:"scopes_supported,omitempty"`

	// BearerMethodsSupported lists the methods for sending bearer tokens.
	BearerMethodsSupported []string `json:"bearer_methods_supported,omitempty"`

	// ResourceSigningAlgValuesSupported lists JWS algorithms for resource operations.
	ResourceSigningAlgValuesSupported []string `json:"resource_signing_alg_values_supported,omitempty"`

	// ResourceName is a human-readable name for the resource.
	ResourceName string `json:"resource_name,omitempty"`

	// ResourceDocumentation is a URL to developer documentation.
	ResourceDocumentation string `json:"resource_documentation,omitempty"`

	// ResourcePolicyURI is a URL to the resource's data usage policy.
	ResourcePolicyURI string `json:"resource_policy_uri,omitempty"`

	// ResourceTOSURI is a URL to the resource's terms of service.
	ResourceTOSURI string `json:"resource_tos_uri,omitempty"`

	// TLSClientCertificateBoundAccessTokens indicates mTLS binding support.
	TLSClientCertificateBoundAccessTokens bool `json:"tls_client_certificate_bound_access_tokens,omitempty"`

	// AuthorizationDetailsTypesSupported lists supported authorization detail types.
	AuthorizationDetailsTypesSupported []string `json:"authorization_details_types_supported,omitempty"`

	// DPOPSigningAlgValuesSupported lists DPoP-supported algorithms.
	DPOPSigningAlgValuesSupported []string `json:"dpop_signing_alg_values_supported,omitempty"`

	// DPOPBoundAccessTokensRequired indicates if DPoP tokens are mandatory.
	DPOPBoundAccessTokensRequired bool `json:"dpop_bound_access_tokens_required,omitempty"`
}

// MetadataPath is the well-known path for Protected Resource Metadata.
const MetadataPath = "/.well-known/oauth-protected-resource"

// BuildMetadata creates Protected Resource Metadata from the OAuth config.
func BuildMetadata(cfg *Config) *ProtectedResourceMetadata {
	metadata := &ProtectedResourceMetadata{
		Resource:               cfg.ResourceServerURL,
		AuthorizationServers:   []string{cfg.Issuer},
		ScopesSupported:        cfg.RequiredScopes,
		BearerMethodsSupported: []string{"header"}, // Bearer token in Authorization header
		ResourceName:           cfg.ResourceName,
		ResourceDocumentation:  cfg.ResourceDocumentation,
	}

	// Add default scopes if none specified
	if len(metadata.ScopesSupported) == 0 {
		metadata.ScopesSupported = []string{"mcp:tools"}
	}

	return metadata
}

// MetadataHandler creates an HTTP handler for the Protected Resource Metadata endpoint.
func MetadataHandler(cfg *Config) http.Handler {
	metadata := BuildMetadata(cfg)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", "GET")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "public, max-age=3600") // Cache for 1 hour

		if err := json.NewEncoder(w).Encode(metadata); err != nil {
			http.Error(w, "Failed to encode metadata", http.StatusInternalServerError)
		}
	})
}

// GetMetadataURL returns the full URL for the Protected Resource Metadata endpoint.
func GetMetadataURL(resourceServerURL string) string {
	// Remove trailing slash if present
	base := strings.TrimSuffix(resourceServerURL, "/")
	return base + MetadataPath
}
