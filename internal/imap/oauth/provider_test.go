package oauth

import (
	"errors"
	"net/http"
	"testing"

	"golang.org/x/oauth2"
)

func TestProviderByName_Google(t *testing.T) {
	p, err := ProviderByName("google")
	if err != nil {
		t.Fatal(err)
	}
	if p.Endpoint().TokenURL == "" {
		t.Fatal("Google provider missing TokenURL")
	}
	if p.DeviceAuthURL() == "" {
		t.Fatal("Google provider missing DeviceAuthURL")
	}
	got := p.Scopes()
	found := false
	for _, s := range got {
		if s == IMAPScope {
			found = true
		}
	}
	if !found {
		t.Fatalf("scopes %v missing %q", got, IMAPScope)
	}
}

func TestProviderByName_Unknown(t *testing.T) {
	if _, err := ProviderByName("yahoo"); err == nil {
		t.Fatal("expected error for unknown provider")
	}
}

func TestIsTerminalAuthError_NilSafe(t *testing.T) {
	if IsTerminalAuthError(nil) {
		t.Fatal("nil should not be terminal")
	}
}

func TestIsTerminalAuthError_InvalidGrant(t *testing.T) {
	rerr := &oauth2.RetrieveError{ErrorCode: "invalid_grant"}
	if !IsTerminalAuthError(rerr) {
		t.Fatal("invalid_grant should be terminal")
	}
}

func TestIsTerminalAuthError_ServerError(t *testing.T) {
	rerr := &oauth2.RetrieveError{Response: &http.Response{StatusCode: 503}}
	if IsTerminalAuthError(rerr) {
		t.Fatal("5xx should not be terminal")
	}
}

func TestIsTerminalAuthError_4xxIsTerminal(t *testing.T) {
	rerr := &oauth2.RetrieveError{Response: &http.Response{StatusCode: 400}}
	if !IsTerminalAuthError(rerr) {
		t.Fatal("4xx should be terminal")
	}
}

func TestIsTerminalAuthError_GenericError(t *testing.T) {
	if IsTerminalAuthError(errors.New("connection reset by peer")) {
		t.Fatal("network errors should not be terminal")
	}
}

func TestIsTerminalAuthError_StringMatch(t *testing.T) {
	if !IsTerminalAuthError(errors.New("oauth2: cannot fetch token: 400 Bad Request\nResponse: {\"error\": \"invalid_grant\"}")) {
		t.Fatal("invalid_grant in message should be terminal")
	}
}
