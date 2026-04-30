package config

import (
	"errors"
	"path/filepath"
	"testing"
)

func TestValidate_PasswordAuthBackwardsCompat(t *testing.T) {
	// A config with no auth block and a password should still validate —
	// existing deployments must keep working unchanged.
	cfg := &Config{
		IMAP: IMAPConfig{
			Host:     "imap.gmail.com",
			Username: "user@example.com",
			Password: "secret",
		},
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestValidate_PasswordMissing(t *testing.T) {
	cfg := &Config{
		IMAP: IMAPConfig{
			Host:     "imap.gmail.com",
			Username: "user@example.com",
		},
	}
	if err := cfg.Validate(); !errors.Is(err, ErrMissingIMAPPassword) {
		t.Fatalf("want ErrMissingIMAPPassword, got %v", err)
	}
}

func TestValidate_XOAUTH2Valid(t *testing.T) {
	cfg := &Config{
		IMAP: IMAPConfig{
			Host:     "imap.gmail.com",
			Username: "user@example.com",
			Auth: IMAPAuthConfig{
				Type:         AuthTypeXOAUTH2,
				Provider:     OAuthProviderGoogle,
				ClientID:     "abc",
				ClientSecret: "xyz",
			},
		},
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestValidate_XOAUTH2MissingClientID(t *testing.T) {
	cfg := &Config{
		IMAP: IMAPConfig{
			Host:     "imap.gmail.com",
			Username: "user@example.com",
			Auth: IMAPAuthConfig{
				Type:         AuthTypeXOAUTH2,
				Provider:     OAuthProviderGoogle,
				ClientSecret: "xyz",
			},
		},
	}
	if err := cfg.Validate(); !errors.Is(err, ErrMissingOAuthClientID) {
		t.Fatalf("want ErrMissingOAuthClientID, got %v", err)
	}
}

func TestValidate_XOAUTH2UnknownProvider(t *testing.T) {
	cfg := &Config{
		IMAP: IMAPConfig{
			Host:     "imap.gmail.com",
			Username: "user@example.com",
			Auth: IMAPAuthConfig{
				Type:         AuthTypeXOAUTH2,
				Provider:     "yahoo",
				ClientID:     "abc",
				ClientSecret: "xyz",
			},
		},
	}
	if err := cfg.Validate(); !errors.Is(err, ErrUnknownOAuthProvider) {
		t.Fatalf("want ErrUnknownOAuthProvider, got %v", err)
	}
}

func TestValidate_UnknownAuthType(t *testing.T) {
	cfg := &Config{
		IMAP: IMAPConfig{
			Host:     "imap.gmail.com",
			Username: "user@example.com",
			Auth:     IMAPAuthConfig{Type: "magic"},
		},
	}
	if err := cfg.Validate(); !errors.Is(err, ErrUnknownAuthType) {
		t.Fatalf("want ErrUnknownAuthType, got %v", err)
	}
}

func TestSecretsPath_CoLocatesWithDB(t *testing.T) {
	cfg := &Config{Database: DatabaseConfig{Path: "/var/lib/parse-dmarc/db.sqlite"}}
	want := filepath.Join("/var/lib/parse-dmarc", "secrets.json")
	if got := cfg.SecretsPath(); got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestIMAPAuthConfig_IsXOAUTH2(t *testing.T) {
	if (IMAPAuthConfig{}).IsXOAUTH2() {
		t.Fatal("empty auth block should not be XOAUTH2")
	}
	if (IMAPAuthConfig{Type: AuthTypePassword}).IsXOAUTH2() {
		t.Fatal("password auth should not be XOAUTH2")
	}
	if !(IMAPAuthConfig{Type: AuthTypeXOAUTH2}).IsXOAUTH2() {
		t.Fatal("xoauth2 auth should be XOAUTH2")
	}
}
