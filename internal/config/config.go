package config

import (
	"errors"
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
)

// Config holds the application configuration
type Config struct {
	IMAP     IMAPConfig     `json:"imap"`
	Database DatabaseConfig `json:"database"`
	Server   ServerConfig   `json:"server"`
}

// IMAPConfig holds IMAP server configuration
type IMAPConfig struct {
	Host     string `json:"host" env:"IMAP_HOST"`
	Port     int    `json:"port" env:"IMAP_PORT" envDefault:"993"`
	Username string `json:"username" env:"IMAP_USERNAME"`
	Password string `json:"password" env:"IMAP_PASSWORD"`
	Mailbox  string `json:"mailbox" env:"IMAP_MAILBOX" envDefault:"INBOX"`
	UseTLS   bool   `json:"use_tls" env:"IMAP_USE_TLS" envDefault:"true"`
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

	if _, err := os.Stat(path); err == nil {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(data, &cfg); err != nil {
			return nil, err
		}
	}

	if err := env.Parse(&cfg); err != nil {
		return nil, err
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
				return nil, err
			}
			err = ensureDBPathExists(cfg.Database.Path)
			if err != nil {
				return nil, err
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

// GenerateSample creates a sample configuration file
func GenerateSample(path string) error {
	dbPath, err := defaultDBPath()
	if err != nil {
		return err
	}
	sample := Config{
		IMAP: IMAPConfig{
			Host:     "imap.example.com",
			Port:     993,
			Username: "your-email@example.com",
			Password: "your-password",
			Mailbox:  "INBOX",
			UseTLS:   true,
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
		return err
	}

	return os.WriteFile(path, data, 0644)
}
