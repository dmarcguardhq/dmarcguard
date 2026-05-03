package auth

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/goccy/go-json"
)

// GitHub OAuth and API endpoints.
const (
	githubAuthURL   = "https://github.com/login/oauth/authorize"
	githubTokenURL  = "https://github.com/login/oauth/access_token"
	githubUserURL   = "https://api.github.com/user"
	githubEmailsURL = "https://api.github.com/user/emails"
	githubScope     = "read:user user:email"
	githubUserAgent = "dmarcguard"
)

// Identity is what the OAuth callback hands off to the allowlist + signer.
type Identity struct {
	Login string // GitHub username
	Email string // verified primary email
}

// GitHubClient wraps the OAuth dance and the user/emails API calls.
// Construct via NewGitHubClient.
type GitHubClient struct {
	clientID     string
	clientSecret string
	redirectURL  string
	httpClient   *http.Client
}

// NewGitHubClient returns a client. http.DefaultClient is used when h is nil.
func NewGitHubClient(clientID, clientSecret, redirectURL string, h *http.Client) *GitHubClient {
	if h == nil {
		h = http.DefaultClient
	}
	return &GitHubClient{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURL:  redirectURL,
		httpClient:   h,
	}
}

// AuthCodeURL builds the GitHub authorization URL for the given state.
func (g *GitHubClient) AuthCodeURL(state string) string {
	q := url.Values{}
	q.Set("client_id", g.clientID)
	q.Set("redirect_uri", g.redirectURL)
	q.Set("scope", githubScope)
	q.Set("state", state)
	q.Set("allow_signup", "false")
	return githubAuthURL + "?" + q.Encode()
}

// Exchange swaps an authorization code for an access token, then fetches the
// authenticated user's verified primary email and login. The returned Identity
// is what the allowlist is checked against.
func (g *GitHubClient) Exchange(ctx context.Context, code string) (*Identity, error) {
	tok, err := g.fetchToken(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("exchange code: %w", err)
	}

	login, err := g.fetchLogin(ctx, tok)
	if err != nil {
		return nil, fmt.Errorf("fetch user: %w", err)
	}

	email, err := g.fetchPrimaryVerifiedEmail(ctx, tok)
	if err != nil {
		return nil, fmt.Errorf("fetch email: %w", err)
	}

	return &Identity{Login: login, Email: email}, nil
}

func (g *GitHubClient) fetchToken(ctx context.Context, code string) (string, error) {
	body := url.Values{}
	body.Set("client_id", g.clientID)
	body.Set("client_secret", g.clientSecret)
	body.Set("code", code)
	body.Set("redirect_uri", g.redirectURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, githubTokenURL, strings.NewReader(body.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", githubUserAgent)

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("github token endpoint returned %d: %s", resp.StatusCode, string(raw))
	}

	var tokResp struct {
		AccessToken string `json:"access_token"`
		Error       string `json:"error"`
		ErrorDesc   string `json:"error_description"`
	}
	if err := json.Unmarshal(raw, &tokResp); err != nil {
		return "", fmt.Errorf("parse token response: %w", err)
	}
	if tokResp.Error != "" {
		return "", fmt.Errorf("github oauth error %q: %s", tokResp.Error, tokResp.ErrorDesc)
	}
	if tokResp.AccessToken == "" {
		return "", errors.New("github returned empty access token")
	}
	return tokResp.AccessToken, nil
}

func (g *GitHubClient) fetchLogin(ctx context.Context, token string) (string, error) {
	var u struct {
		Login string `json:"login"`
	}
	if err := g.getJSON(ctx, githubUserURL, token, &u); err != nil {
		return "", err
	}
	if u.Login == "" {
		return "", errors.New("github user response missing login")
	}
	return u.Login, nil
}

func (g *GitHubClient) fetchPrimaryVerifiedEmail(ctx context.Context, token string) (string, error) {
	var emails []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}
	if err := g.getJSON(ctx, githubEmailsURL, token, &emails); err != nil {
		return "", err
	}
	for _, e := range emails {
		if e.Primary && e.Verified {
			return e.Email, nil
		}
	}
	// Fallback: any verified email.
	for _, e := range emails {
		if e.Verified {
			return e.Email, nil
		}
	}
	return "", errors.New("no verified email on github account")
}

func (g *GitHubClient) getJSON(ctx context.Context, urlStr, token string, into interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", githubUserAgent)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("github api %s returned %d: %s", urlStr, resp.StatusCode, string(raw))
	}
	return json.NewDecoder(resp.Body).Decode(into)
}
