package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAllowlist_EmailMatch(t *testing.T) {
	a := NewAllowlist([]string{"Alice@Example.com"}, nil)
	if !a.Allows(&Identity{Login: "alice", Email: "alice@example.com"}) {
		t.Fatal("email should match case-insensitively")
	}
	if a.Allows(&Identity{Login: "bob", Email: "bob@example.com"}) {
		t.Fatal("non-allowed email should be rejected")
	}
}

func TestAllowlist_LoginMatch(t *testing.T) {
	a := NewAllowlist(nil, []string{"sebykrueger"})
	if !a.Allows(&Identity{Login: "SebyKrueger", Email: "x@example.com"}) {
		t.Fatal("login should match case-insensitively")
	}
}

func TestAllowlist_EmptyDeniesAll(t *testing.T) {
	a := NewAllowlist(nil, nil)
	if a.Allows(&Identity{Login: "anyone", Email: "anyone@example.com"}) {
		t.Fatal("empty allowlist must deny everyone")
	}
}

func TestMiddleware_ProtectedRedirectsBrowser(t *testing.T) {
	signer := newSigner(t)
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mw := Middleware(signer, next)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	mw.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Fatalf("expected 303, got %d", rr.Code)
	}
	if loc := rr.Header().Get("Location"); loc != "/auth/login" {
		t.Fatalf("expected Location=/auth/login, got %q", loc)
	}
}

func TestMiddleware_ApiReturns401Json(t *testing.T) {
	signer := newSigner(t)
	mw := Middleware(signer, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {}))

	req := httptest.NewRequest(http.MethodGet, "/api/statistics", nil)
	rr := httptest.NewRecorder()
	mw.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected JSON, got %q", ct)
	}
}

func TestMiddleware_ValidSessionPassesThrough(t *testing.T) {
	signer := newSigner(t)
	cookie, _ := signer.Sign("alice", "a@x.com", time.Now())

	called := false
	mw := Middleware(signer, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		called = true
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/statistics", nil)
	req.AddCookie(&http.Cookie{Name: CookieName, Value: cookie})
	rr := httptest.NewRecorder()
	mw.ServeHTTP(rr, req)

	if !called {
		t.Fatal("next handler should have been invoked for valid session")
	}
}
