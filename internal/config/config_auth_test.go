package config

import (
	"errors"
	"strings"
	"testing"
)

func validAuth() AuthConfig {
	return AuthConfig{
		Enabled:       true,
		ClientID:      "abc",
		ClientSecret:  "xyz",
		RedirectURL:   "https://dmarc.example.com/auth/callback",
		SessionSecret: strings.Repeat("a", 44),
		AllowedEmails: []string{"seb@example.com"},
	}
}

func TestValidateAuth_DisabledIsAlwaysOK(t *testing.T) {
	cfg := &Config{Auth: AuthConfig{Enabled: false}}
	if err := cfg.ValidateAuth(); err != nil {
		t.Fatalf("disabled auth should not error: %v", err)
	}
}

func TestValidateAuth_HappyPath(t *testing.T) {
	cfg := &Config{Auth: validAuth()}
	if err := cfg.ValidateAuth(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestValidateAuth_MissingClientID(t *testing.T) {
	a := validAuth()
	a.ClientID = ""
	cfg := &Config{Auth: a}
	if err := cfg.ValidateAuth(); !errors.Is(err, ErrAuthMissingClientID) {
		t.Fatalf("want ErrAuthMissingClientID, got %v", err)
	}
}

func TestValidateAuth_MissingSecret(t *testing.T) {
	a := validAuth()
	a.ClientSecret = ""
	cfg := &Config{Auth: a}
	if err := cfg.ValidateAuth(); !errors.Is(err, ErrAuthMissingClientSecret) {
		t.Fatalf("want ErrAuthMissingClientSecret, got %v", err)
	}
}

func TestValidateAuth_ShortSessionSecret(t *testing.T) {
	a := validAuth()
	a.SessionSecret = "tooshort"
	cfg := &Config{Auth: a}
	if err := cfg.ValidateAuth(); !errors.Is(err, ErrAuthSessionSecretTooShort) {
		t.Fatalf("want ErrAuthSessionSecretTooShort, got %v", err)
	}
}

func TestValidateAuth_EmptyAllowlistRejected(t *testing.T) {
	a := validAuth()
	a.AllowedEmails = nil
	a.AllowedUsers = nil
	cfg := &Config{Auth: a}
	if err := cfg.ValidateAuth(); !errors.Is(err, ErrAuthAllowlistEmpty) {
		t.Fatalf("want ErrAuthAllowlistEmpty, got %v", err)
	}
}

func TestValidateAuth_AllowsUserOnlyAllowlist(t *testing.T) {
	a := validAuth()
	a.AllowedEmails = nil
	a.AllowedUsers = []string{"sebykrueger"}
	cfg := &Config{Auth: a}
	if err := cfg.ValidateAuth(); err != nil {
		t.Fatalf("user-only allowlist should validate, got %v", err)
	}
}
