package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/caarlos0/env/v11"
	"github.com/goccy/go-json"
)

var (
	// ErrMissingIMAPHost is returned when IMAP host is not configured
	ErrMissingIMAPHost = errors.New("IMAP_HOST is required: set via environment variable or config file")
	// ErrMissingIMAPUsername is returned when IMAP username is not configured
	ErrMissingIMAPUsername = errors.New("IMAP_USERNAME is required: set via environment variable or config file")
	// ErrMissingIMAPPassword is returned when IMAP password is not configured
	ErrMissingIMAPPassword = errors.New("IMAP_PASSWORD is required: set via environment variable or config file")
	// ErrAuthMissingClientID is returned when auth is enabled without a GitHub OAuth client ID
	ErrAuthMissingClientID = errors.New("auth.client_id is required when auth.enabled (set AUTH_CLIENT_ID)")
	// ErrAuthMissingClientSecret is returned when auth is enabled without a GitHub OAuth client secret
	ErrAuthMissingClientSecret = errors.New("auth.client_secret is required when auth.enabled (set AUTH_CLIENT_SECRET)")
	// ErrAuthMissingRedirectURL is returned when auth is enabled without a callback URL
	ErrAuthMissingRedirectURL = errors.New("auth.redirect_url is required when auth.enabled (e.g., https://dmarc.example.com/auth/callback)")
	// ErrAuthMissingSessionSecret is returned when auth is enabled without a session-signing secret
	ErrAuthMissingSessionSecret = errors.New("auth.session_secret is required when auth.enabled (generate with --gen-session-secret)")
	// ErrAuthSessionSecretTooShort is returned when the session secret is too short to safely sign cookies
	ErrAuthSessionSecretTooShort = errors.New("auth.session_secret must decode to at least 32 bytes")
	// ErrAuthAllowlistEmpty is returned when auth is enabled but no users/emails are allowed
	ErrAuthAllowlistEmpty = errors.New("auth requires at least one entry in allowed_emails or allowed_users (otherwise no one can log in)")
)

// Config holds the application configuration
type Config struct {
	LogLevel    string         `json:"log_level" env:"LOG_LEVEL" envDefault:"info"`
	ColoredLogs bool           `json:"colored_logs" env:"COLORED_LOGS" envDefault:"false"`
	IMAP        IMAPConfig     `json:"imap"`
	Database    DatabaseConfig `json:"database"`
	Server      ServerConfig   `json:"server"`
	Auth        AuthConfig     `json:"auth" envPrefix:"AUTH_"`
}

// AuthConfig holds dashboard authentication configuration. When Enabled is
// false (or the block is absent) the dashboard runs unauthenticated, matching
// pre-auth behavior. When enabled, all /api/* and / routes require a valid
// session cookie obtained via GitHub OAuth.
type AuthConfig struct {
	Enabled        bool     `json:"enabled" env:"ENABLED"`
	ClientID       string   `json:"client_id" env:"CLIENT_ID"`
	ClientSecret   string   `json:"client_secret" env:"CLIENT_SECRET"`
	RedirectURL    string   `json:"redirect_url" env:"REDIRECT_URL"`
	SessionSecret  string   `json:"session_secret" env:"SESSION_SECRET"`
	AllowedEmails  []string `json:"allowed_emails" env:"ALLOWED_EMAILS" envSeparator:","`
	AllowedUsers   []string `json:"allowed_users" env:"ALLOWED_USERS" envSeparator:","`
	SessionTTLDays int      `json:"session_ttl_days" env:"SESSION_TTL_DAYS"`
}

// IMAPConfig holds IMAP server configuration
type IMAPConfig struct {
	Host     string `json:"host" env:"IMAP_HOST"`
	Port     int    `json:"port" env:"IMAP_PORT" envDefault:"993"`
	Username string `json:"username" env:"IMAP_USERNAME"`
	Password string `json:"password" env:"IMAP_PASSWORD"`
	Mailbox  string `json:"mailbox" env:"IMAP_MAILBOX" envDefault:"INBOX"`
	UseTLS   bool   `json:"use_tls" env:"IMAP_USE_TLS" envDefault:"true"`

	MarkAsSeen       bool   `json:"mark_as_seen" env:"IMAP_MARK_AS_SEEN" envDefault:"true"`
	ProcessedMailbox string `json:"processed_mailbox" env:"IMAP_PROCESSED_MAILBOX"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Path string `json:"path" env:"DATABASE_PATH"`
}

// ServerConfig holds web server configuration
type ServerConfig struct {
	Port int    `json:"port" env:"SERVER_PORT" envDefault:"8080"`
	Host string `json:"host" env:"SERVER_HOST" envDefault:""`
}

func defaultDBPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return "", errors.New("cannot determine home directory")
	}
	return filepath.Join(home, ".parse-dmarc/db.sqlite"), nil
}

func fallbackDBPath() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", errors.New("cannot determine home directory or current working directory")
	}
	return filepath.Join(cwd, ".parse-dmarc/db.sqlite"), nil
}

func ensureDBPathExists(dbPath string) error {
	parent := filepath.Dir(dbPath)
	if err := os.MkdirAll(parent, 0755); err != nil {
		return errors.New("failed to create database directory at " + parent + ": " + err.Error() + " - ensure the path is writable or set DATABASE_PATH environment variable")
	}
	return nil
}

// Load loads configuration from a JSON file
func Load(path string) (*Config, error) {
	var cfg Config
	var err error

	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("parse env config: %w", err)
	}

	if _, err := os.Stat(path); err == nil {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read config file %s: %w", path, err)
		}

		if err := json.Unmarshal(data, &cfg); err != nil {
			return nil, fmt.Errorf("parse config file %s: %w", path, err)
		}
	}

	if cfg.IMAP.Port == 0 {
		cfg.IMAP.Port = 993
	}
	if cfg.IMAP.Mailbox == "" {
		cfg.IMAP.Mailbox = "INBOX"
	}
	if cfg.Database.Path == "" {
		cfg.Database.Path, err = defaultDBPath()
		if err != nil || ensureDBPathExists(cfg.Database.Path) != nil {
			cfg.Database.Path, err = fallbackDBPath()
			if err != nil {
				return nil, fmt.Errorf("resolve database path: %w", err)
			}
			err = ensureDBPathExists(cfg.Database.Path)
			if err != nil {
				return nil, fmt.Errorf("ensure database path: %w", err)
			}
		}
	}
	if cfg.Server.Port == 0 {
		cfg.Server.Port = 8080
	}

	return &cfg, nil
}

// Validate checks that all required configuration values are set.
// Required fields: IMAP host, username, and password.
// Returns nil if valid, or an error describing the missing configuration.
func (c *Config) Validate() error {
	if c.IMAP.Host == "" {
		return ErrMissingIMAPHost
	}
	if c.IMAP.Username == "" {
		return ErrMissingIMAPUsername
	}
	if c.IMAP.Password == "" {
		return ErrMissingIMAPPassword
	}
	return nil
}

// ValidateAuth checks the auth configuration when it is enabled. Always called
// before mounting auth handlers. Returns nil when auth is disabled.
func (c *Config) ValidateAuth() error {
	if !c.Auth.Enabled {
		return nil
	}
	if c.Auth.ClientID == "" {
		return ErrAuthMissingClientID
	}
	if c.Auth.ClientSecret == "" {
		return ErrAuthMissingClientSecret
	}
	if c.Auth.RedirectURL == "" {
		return ErrAuthMissingRedirectURL
	}
	if c.Auth.SessionSecret == "" {
		return ErrAuthMissingSessionSecret
	}
	// Decoded length check happens when the secret is loaded by the auth
	// package; here we just guard against obviously-short raw strings.
	if len(c.Auth.SessionSecret) < 32 {
		return ErrAuthSessionSecretTooShort
	}
	if len(c.Auth.AllowedEmails) == 0 && len(c.Auth.AllowedUsers) == 0 {
		return ErrAuthAllowlistEmpty
	}
	return nil
}

// GenerateSample creates a sample configuration file
func GenerateSample(path string) error {
	dbPath, err := defaultDBPath()
	if err != nil {
		return fmt.Errorf("resolve default database path: %w", err)
	}
	sample := Config{
		LogLevel: "info",
		IMAP: IMAPConfig{
			Host:     "imap.example.com",
			Port:     993,
			Username: "your-email@example.com",
			Password: "your-password",
			Mailbox:  "INBOX",
			UseTLS:   true,

			MarkAsSeen:       true,
			ProcessedMailbox: "",
		},
		Database: DatabaseConfig{
			Path: dbPath,
		},
		Server: ServerConfig{
			Port: 8080,
			Host: "0.0.0.0",
		},
	}

	data, err := json.MarshalIndent(sample, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal sample config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write config file %s: %w", path, err)
	}

	return nil
}
