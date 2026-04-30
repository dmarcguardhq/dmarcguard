// Package oauth implements OAuth2 token acquisition and storage for IMAP
// XOAUTH2 authentication. It supports the OAuth2 device authorization grant
// (RFC 8628) for headless bootstrap and persists refresh tokens to a local
// secrets file co-located with the database.
package oauth

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-json"
)

// EnvRefreshToken is the env var that, when set, overrides the on-disk
// refresh token. Read-only: the daemon will never write back to it, even if
// the provider rotates the refresh token.
const EnvRefreshToken = "IMAP_OAUTH_REFRESH_TOKEN"

// ErrNoRefreshToken indicates no refresh token is available — the user must
// run --oauth-login before the daemon can authenticate.
var ErrNoRefreshToken = errors.New("no refresh token available: run with --oauth-login to authorize")

// Secrets is the on-disk shape of secrets.json.
type Secrets struct {
	RefreshToken string `json:"refresh_token"`
}

// LoadSecrets reads the refresh token, preferring the env override and
// falling back to the secrets file. Returns ErrNoRefreshToken when neither
// source has a token. envOverride indicates whether the token came from env
// (and therefore must not be written back).
func LoadSecrets(path string) (token string, envOverride bool, err error) {
	if v := os.Getenv(EnvRefreshToken); v != "" {
		return v, true, nil
	}

	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return "", false, ErrNoRefreshToken
	}
	if err != nil {
		return "", false, fmt.Errorf("read secrets file %s: %w", path, err)
	}

	var s Secrets
	if err := json.Unmarshal(data, &s); err != nil {
		return "", false, fmt.Errorf("parse secrets file %s: %w", path, err)
	}
	if s.RefreshToken == "" {
		return "", false, ErrNoRefreshToken
	}
	return s.RefreshToken, false, nil
}

// SaveSecrets writes the refresh token atomically (write-tmp + rename) with
// 0600 permissions, creating parent dirs as needed.
func SaveSecrets(path, refreshToken string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create secrets directory: %w", err)
	}

	data, err := json.MarshalIndent(Secrets{RefreshToken: refreshToken}, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal secrets: %w", err)
	}

	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o600); err != nil {
		return fmt.Errorf("write secrets file: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("rename secrets file: %w", err)
	}
	return nil
}
