package auth

import (
	"context"
	"net/http"
	"strings"
	"time"
)

// contextKey is a private type for context keys to avoid collisions
type contextKey string

const userContextKey contextKey = "user"

// Allowlist holds the set of GitHub usernames and verified emails authorized
// to access the dashboard. Matching is case-insensitive on both sides.
type Allowlist struct {
	emails map[string]struct{}
	users  map[string]struct{}
}

// NewAllowlist normalizes and indexes the configured allowlist entries.
func NewAllowlist(emails, users []string) *Allowlist {
	a := &Allowlist{
		emails: make(map[string]struct{}, len(emails)),
		users:  make(map[string]struct{}, len(users)),
	}
	for _, e := range emails {
		if e = strings.TrimSpace(strings.ToLower(e)); e != "" {
			a.emails[e] = struct{}{}
		}
	}
	for _, u := range users {
		if u = strings.TrimSpace(strings.ToLower(u)); u != "" {
			a.users[u] = struct{}{}
		}
	}
	return a
}

// Allows returns true if the identity matches at least one allowlist entry.
func (a *Allowlist) Allows(id *Identity) bool {
	if _, ok := a.emails[strings.ToLower(id.Email)]; ok {
		return true
	}
	if _, ok := a.users[strings.ToLower(id.Login)]; ok {
		return true
	}
	return false
}

// UserFromContext extracts the Identity from the request context.
// Returns nil if no identity is present (e.g., unauthenticated request).
func UserFromContext(ctx context.Context) *Identity {
	if id, ok := ctx.Value(userContextKey).(*Identity); ok {
		return id
	}
	return nil
}

// Middleware wraps next with a session check. Unauthenticated requests are
// either redirected to /auth/login (for HTML/dashboard requests) or get a
// 401 JSON response (for /api/* requests). The split avoids breaking
// programmatic API consumers that should see a 401 rather than a 302.
func Middleware(signer *SessionSigner, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(CookieName)
		if err != nil {
			deny(w, r)
			return
		}
		session, err := signer.Verify(cookie.Value, time.Now())
		if err != nil {
			deny(w, r)
			return
		}
		// Reconstruct the Identity from the session (Subject was set to the
		// GitHub login at sign-time). UserFromContext expects *Identity.
		identity := &Identity{Login: session.Subject, Email: session.Email}
		ctx := context.WithValue(r.Context(), userContextKey, identity)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func deny(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/api/") {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"unauthorized","login_url":"/auth/login"}`))
		return
	}
	http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
}
