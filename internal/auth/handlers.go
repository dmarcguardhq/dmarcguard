package auth

import (
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

// Handlers carries the dependencies for the /auth/* HTTP handlers.
type Handlers struct {
	GitHub    *GitHubClient
	Signer    *SessionSigner
	Allowlist *Allowlist
	Log       *zerolog.Logger
	// Secure controls whether session/state cookies set the Secure flag.
	// Disabled when the redirect URL is http:// (typically local dev).
	Secure bool
}

// NewHandlers constructs Handlers and infers the Secure flag from the
// redirect URL scheme.
func NewHandlers(github *GitHubClient, signer *SessionSigner, allowlist *Allowlist, redirectURL string, log *zerolog.Logger) *Handlers {
	return &Handlers{
		GitHub:    github,
		Signer:    signer,
		Allowlist: allowlist,
		Log:       log,
		Secure:    strings.HasPrefix(redirectURL, "https://"),
	}
}

// Mount attaches the /auth/* handlers to mux.
func (h *Handlers) Mount(mux *http.ServeMux) {
	mux.HandleFunc("/auth/login", h.Login)
	mux.HandleFunc("/auth/callback", h.Callback)
	mux.HandleFunc("/auth/logout", h.Logout)
	mux.HandleFunc("/auth/logged-out", h.LoggedOut)
}

// Login generates a fresh state token, sets it as a short-lived cookie, and
// redirects to GitHub's authorization endpoint.
func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	state, err := h.Signer.SignState(time.Now())
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     StateCookieName,
		Value:    state,
		Path:     "/auth/callback",
		MaxAge:   600,
		HttpOnly: true,
		Secure:   h.Secure,
		SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, h.GitHub.AuthCodeURL(state), http.StatusSeeOther)
}

// Callback validates the OAuth state, exchanges the code for an identity,
// checks the allowlist, signs a session cookie, and redirects to /.
func (h *Handlers) Callback(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	if e := q.Get("error"); e != "" {
		h.Log.Warn().Str("error", e).Str("desc", q.Get("error_description")).Msg("github oauth error")
		http.Error(w, "GitHub authorization failed: "+e, http.StatusBadRequest)
		return
	}
	code := q.Get("code")
	state := q.Get("state")
	if code == "" || state == "" {
		http.Error(w, "missing code or state", http.StatusBadRequest)
		return
	}

	stateCookie, err := r.Cookie(StateCookieName)
	if err != nil || stateCookie.Value != state {
		http.Error(w, "invalid OAuth state (possible CSRF or stale link)", http.StatusBadRequest)
		return
	}
	if err := h.Signer.VerifyState(state, time.Now()); err != nil {
		http.Error(w, "OAuth state expired or tampered", http.StatusBadRequest)
		return
	}

	// Clear the state cookie immediately — single-use.
	http.SetCookie(w, &http.Cookie{
		Name:     StateCookieName,
		Value:    "",
		Path:     "/auth/callback",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   h.Secure,
		SameSite: http.SameSiteLaxMode,
	})

	id, err := h.GitHub.Exchange(r.Context(), code)
	if err != nil {
		h.Log.Warn().Err(err).Msg("github exchange failed")
		http.Error(w, "GitHub authentication failed", http.StatusBadGateway)
		return
	}

	if !h.Allowlist.Allows(id) {
		h.Log.Warn().Str("login", id.Login).Str("email", id.Email).Msg("login denied: not on allowlist")
		http.Error(w, "Account "+id.Login+" ("+id.Email+") is not authorized to access this dashboard.", http.StatusForbidden)
		return
	}

	cookieValue, err := h.Signer.Sign(id.Login, id.Email, time.Now())
	if err != nil {
		h.Log.Error().Err(err).Msg("session sign failed")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    cookieValue,
		Path:     "/",
		MaxAge:   int(h.Signer.TTL().Seconds()),
		HttpOnly: true,
		Secure:   h.Secure,
		SameSite: http.SameSiteLaxMode,
	})

	h.Log.Info().Str("login", id.Login).Str("email", id.Email).Msg("login succeeded")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Logout clears the session cookie and redirects to a dedicated "signed out"
// page. We deliberately do NOT redirect to /auth/login here — that would
// immediately bounce the user back through GitHub, and (because GitHub
// remembers the prior consent) they'd be silently re-authenticated and end
// up exactly where they started. The dedicated landing page makes the user
// click an explicit "Sign in" link to come back.
//
// Accepts both POST (for forms) and GET (for plain "logout" links) — the
// cookie is HttpOnly and SameSite=Lax, so cross-site triggering still
// requires user navigation.
func (h *Handlers) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   h.Secure,
		SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, "/auth/logged-out", http.StatusSeeOther)
}

// LoggedOut is the post-logout landing page. Plain HTML (no Vue) so the
// browser can render it without re-authenticating against the gated SPA.
func (h *Handlers) LoggedOut(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	_, _ = w.Write([]byte(`<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>Signed out — dmarcguard</title>
<style>
  :root { color-scheme: light dark; }
  body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", system-ui, sans-serif; max-width: 28rem; margin: 6rem auto; padding: 2rem; line-height: 1.6; }
  h1 { font-size: 1.5rem; margin: 0 0 0.5rem; }
  p { opacity: 0.75; margin: 0 0 1.5rem; }
  a.btn { display: inline-block; padding: 0.625rem 1.25rem; background: #1f2937; color: #fff; text-decoration: none; border-radius: 0.375rem; font-weight: 500; }
  a.btn:hover { background: #111827; }
  @media (prefers-color-scheme: dark) {
    body { background: #0f172a; color: #e2e8f0; }
    a.btn { background: #e2e8f0; color: #0f172a; }
    a.btn:hover { background: #fff; }
  }
</style>
</head>
<body>
  <h1>Signed out</h1>
  <p>You've been signed out of dmarcguard. Your session cookie has been cleared.</p>
  <p><a class="btn" href="/auth/login">Sign in again</a></p>
</body>
</html>`))
}
