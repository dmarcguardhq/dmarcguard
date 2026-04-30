package oauth

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadSecrets_FileMissing(t *testing.T) {
	t.Setenv(EnvRefreshToken, "")
	_, _, err := LoadSecrets(filepath.Join(t.TempDir(), "nope.json"))
	if !errors.Is(err, ErrNoRefreshToken) {
		t.Fatalf("want ErrNoRefreshToken, got %v", err)
	}
}

func TestSaveLoadSecrets_RoundTrip(t *testing.T) {
	t.Setenv(EnvRefreshToken, "")
	path := filepath.Join(t.TempDir(), "nested", "secrets.json")
	if err := SaveSecrets(path, "rt-12345"); err != nil {
		t.Fatalf("save: %v", err)
	}
	tok, env, err := LoadSecrets(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if env {
		t.Fatal("file-loaded token should not be marked env override")
	}
	if tok != "rt-12345" {
		t.Fatalf("got %q want rt-12345", tok)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if perm := info.Mode().Perm(); perm != 0o600 {
		t.Fatalf("secrets file should be 0600, got %o", perm)
	}
}

func TestLoadSecrets_EnvOverride(t *testing.T) {
	path := filepath.Join(t.TempDir(), "secrets.json")
	if err := SaveSecrets(path, "from-file"); err != nil {
		t.Fatal(err)
	}
	t.Setenv(EnvRefreshToken, "from-env")
	tok, env, err := LoadSecrets(path)
	if err != nil {
		t.Fatal(err)
	}
	if !env {
		t.Fatal("expected envOverride=true")
	}
	if tok != "from-env" {
		t.Fatalf("env should win, got %q", tok)
	}
}

func TestLoadSecrets_EmptyFileReturnsErrNoRefreshToken(t *testing.T) {
	t.Setenv(EnvRefreshToken, "")
	path := filepath.Join(t.TempDir(), "secrets.json")
	if err := os.WriteFile(path, []byte(`{}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, _, err := LoadSecrets(path); !errors.Is(err, ErrNoRefreshToken) {
		t.Fatalf("want ErrNoRefreshToken, got %v", err)
	}
}
