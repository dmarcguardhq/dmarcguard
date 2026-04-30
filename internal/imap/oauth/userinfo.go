package oauth

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/goccy/go-json"
)

// LoginResult is what --oauth-login produces: a refresh token plus the
// authenticated identity, so the user knows which mailbox they just authorized.
type LoginResult struct {
	RefreshToken string
	Email        string
}

// PromptFunc is invoked once the OAuth flow has a URL the user must visit.
// For loopback flow, the second argument is unused (no separate user code).
type PromptFunc func(authURL, userCode string)

// DeviceLoginResult is retained as an alias of LoginResult so callers built
// against the original device-flow API keep compiling. New code should use
// LoginResult directly.
type DeviceLoginResult = LoginResult

func fetchUserinfoEmail(ctx context.Context, p Provider, accessToken string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.UserinfoURL(), nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("userinfo: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var info struct {
		Email string `json:"email"`
	}
	if err := json.Unmarshal(body, &info); err != nil {
		return "", err
	}
	return info.Email, nil
}
