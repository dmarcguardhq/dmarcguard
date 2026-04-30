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
	// ErrMissingOAuthClientID is returned when an XOAUTH2 auth block lacks a client ID
	ErrMissingOAuthClientID = errors.New("imap.auth.client_id is required for xoauth2 (set IMAP_AUTH_CLIENT_ID)")
	// ErrMissingOAuthClientSecret is returned when an XOAUTH2 auth block lacks a client secret
	ErrMissingOAuthClientSecret = errors.New("imap.auth.client_secret is required for xoauth2 (set IMAP_AUTH_CLIENT_SECRET)")
	// ErrUnknownOAuthProvider is returned when the auth provider is not supported
	ErrUnknownOAuthProvider = errors.New("imap.auth.provider must be one of: google")
	// ErrUnknownAuthType is returned when imap.auth.type is set to an unsupported value
	ErrUnknownAuthType = errors.New(`imap.auth.type must be one of: "password", "xoauth2"`)
)

// Auth method identifiers used in IMAPAuthConfig.Type.
const (
	AuthTypePassword = "password"
	AuthTypeXOAUTH2  = "xoauth2"
)

// OAuth provider identifiers used in IMAPAuthConfig.Provider.
const (
	OAuthProviderGoogle = "google"
)

// Config holds the application configuration
type Config struct {
	LogLevel    string         `json:"log_level" env:"LOG_LEVEL" envDefault:"info"`
	ColoredLogs bool           `json:"colored_logs" env:"COLORED_LOGS" envDefault:"false"`
	IMAP        IMAPConfig     `json:"imap"`
	Database    DatabaseConfig `json:"database"`
	Server      ServerConfig   `json:"server"`
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

	Auth IMAPAuthConfig `json:"auth" envPrefix:"IMAP_AUTH_"`
}

// IMAPAuthConfig selects how the IMAP client authenticates.
// When Type is empty the daemon falls back to password auth (backwards compat).
type IMAPAuthConfig struct {
	Type         string `json:"type" env:"TYPE"`
	Provider     string `json:"provider" env:"PROVIDER"`
	ClientID     string `json:"client_id" env:"CLIENT_ID"`
	ClientSecret string `json:"client_secret" env:"CLIENT_SECRET"`
}

// IsXOAUTH2 reports whether the auth block selects XOAUTH2.
func (a IMAPAuthConfig) IsXOAUTH2() bool {
	return a.Type == AuthTypeXOAUTH2
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
// Host and username are always required. Password is required for password auth;
// client_id/client_secret are required for xoauth2 auth.
func (c *Config) Validate() error {
	if c.IMAP.Host == "" {
		return ErrMissingIMAPHost
	}
	if c.IMAP.Username == "" {
		return ErrMissingIMAPUsername
	}

	switch c.IMAP.Auth.Type {
	case "", AuthTypePassword:
		if c.IMAP.Password == "" {
			return ErrMissingIMAPPassword
		}
	case AuthTypeXOAUTH2:
		if c.IMAP.Auth.Provider != OAuthProviderGoogle {
			return ErrUnknownOAuthProvider
		}
		if c.IMAP.Auth.ClientID == "" {
			return ErrMissingOAuthClientID
		}
		if c.IMAP.Auth.ClientSecret == "" {
			return ErrMissingOAuthClientSecret
		}
	default:
		return ErrUnknownAuthType
	}
	return nil
}

// SecretsPath returns the on-disk path for the OAuth refresh token store,
// co-located with the database file (e.g. ~/.parse-dmarc/secrets.json).
func (c *Config) SecretsPath() string {
	return filepath.Join(filepath.Dir(c.Database.Path), "secrets.json")
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
