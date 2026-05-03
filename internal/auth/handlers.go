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

// Mount attaches /auth/login, /auth/callback, and /auth/logout to mux.
func (h *Handlers) Mount(mux *http.ServeMux) {
	mux.HandleFunc("/auth/login", h.Login)
	mux.HandleFunc("/auth/callback", h.Callback)
	mux.HandleFunc("/auth/logout", h.Logout)
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

// Logout clears the session cookie and redirects to /auth/login.
// Accepts both POST (for forms) and GET (for "logout" links) — the latter
// is fine because the cookie is HttpOnly and SameSite=Lax, so cross-site
// triggering still requires user navigation.
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
	http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
}
