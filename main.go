package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/meysam81/parse-dmarc/internal/api"
	"github.com/meysam81/parse-dmarc/internal/config"
	"github.com/meysam81/parse-dmarc/internal/imap"
	imapoauth "github.com/meysam81/parse-dmarc/internal/imap/oauth"
	"github.com/meysam81/parse-dmarc/internal/logger"
	mcpserver "github.com/meysam81/parse-dmarc/internal/mcp"
	"github.com/meysam81/parse-dmarc/internal/mcp/oauth"
	"github.com/meysam81/parse-dmarc/internal/metrics"
	"github.com/meysam81/parse-dmarc/internal/parser"
	"github.com/meysam81/parse-dmarc/internal/storage"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v3"
	"golang.org/x/oauth2"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"

	log *zerolog.Logger
)

func main() {
	cli.VersionPrinter = func(c *cli.Command) {
		fmt.Println(version)
	}
	cmd := &cli.Command{
		Name:                  "parse-dmarc",
		Usage:                 "Parse and analyze DMARC reports from IMAP mailbox",
		Version:               version,
		EnableShellCompletion: true,
		Suggest:               true,
		Before: func(ctx context.Context, c *cli.Command) (context.Context, error) {
			log = logger.NewLogger("info", false)
			return ctx, nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Path to configuration file",
				Value:   "config.json",
				Sources: cli.EnvVars("PARSE_DMARC_CONFIG"),
			},
			&cli.BoolFlag{
				Name:    "gen-config",
				Usage:   "Generate sample configuration file",
				Sources: cli.EnvVars("PARSE_DMARC_GEN_CONFIG"),
			},
			&cli.BoolFlag{
				Name:    "fetch-once",
				Usage:   "Fetch reports once and exit",
				Sources: cli.EnvVars("PARSE_DMARC_FETCH_ONCE"),
			},
			&cli.BoolFlag{
				Name:    "serve-only",
				Usage:   "Only serve the dashboard without fetching",
				Sources: cli.EnvVars("PARSE_DMARC_SERVE_ONLY"),
			},
			&cli.IntFlag{
				Name:    "fetch-interval",
				Usage:   "Interval in seconds between fetch operations",
				Value:   300,
				Sources: cli.EnvVars("PARSE_DMARC_FETCH_INTERVAL"),
			},
			&cli.BoolFlag{
				Name:    "metrics",
				Usage:   "Enable Prometheus metrics endpoint at /metrics",
				Value:   true,
				Sources: cli.EnvVars("PARSE_DMARC_METRICS"),
			},
			&cli.BoolFlag{
				Name:    "mcp",
				Usage:   "Run as MCP (Model Context Protocol) server over stdio",
				Sources: cli.EnvVars("PARSE_DMARC_MCP"),
			},
			&cli.StringFlag{
				Name:    "mcp-http",
				Usage:   "Run MCP server over HTTP/SSE at the specified address (e.g., :8081)",
				Sources: cli.EnvVars("PARSE_DMARC_MCP_HTTP"),
			},
			// OAuth2 flags for MCP HTTP server
			&cli.BoolFlag{
				Name:    "mcp-oauth",
				Usage:   "Enable OAuth2 authentication for MCP HTTP server",
				Sources: cli.EnvVars("PARSE_DMARC_MCP_OAUTH"),
			},
			&cli.StringFlag{
				Name:    "mcp-oauth-issuer",
				Usage:   "OAuth2/OIDC issuer URL (e.g., https://auth.example.com/realms/master)",
				Sources: cli.EnvVars("PARSE_DMARC_MCP_OAUTH_ISSUER"),
			},
			&cli.StringFlag{
				Name:    "mcp-oauth-audience",
				Usage:   "Expected audience claim in tokens (usually the MCP server URL)",
				Sources: cli.EnvVars("PARSE_DMARC_MCP_OAUTH_AUDIENCE"),
			},
			&cli.StringFlag{
				Name:    "mcp-oauth-client-id",
				Usage:   "OAuth2 client ID for token introspection",
				Sources: cli.EnvVars("PARSE_DMARC_MCP_OAUTH_CLIENT_ID"),
			},
			&cli.StringFlag{
				Name:    "mcp-oauth-client-secret",
				Usage:   "OAuth2 client secret for token introspection",
				Sources: cli.EnvVars("PARSE_DMARC_MCP_OAUTH_CLIENT_SECRET"),
			},
			&cli.StringFlag{
				Name:    "mcp-oauth-scopes",
				Usage:   "Required scopes (comma-separated, e.g., mcp:tools)",
				Value:   "mcp:tools",
				Sources: cli.EnvVars("PARSE_DMARC_MCP_OAUTH_SCOPES"),
			},
			&cli.StringFlag{
				Name:    "mcp-oauth-introspection-endpoint",
				Usage:   "Token introspection endpoint URL (if set, uses introspection instead of JWT validation)",
				Sources: cli.EnvVars("PARSE_DMARC_MCP_OAUTH_INTROSPECTION_ENDPOINT"),
			},
			&cli.StringFlag{
				Name:    "mcp-oauth-resource-name",
				Usage:   "Human-readable name for the MCP server (for metadata)",
				Value:   "Parse-DMARC MCP Server",
				Sources: cli.EnvVars("PARSE_DMARC_MCP_OAUTH_RESOURCE_NAME"),
			},
			&cli.BoolFlag{
				Name:    "mcp-oauth-insecure",
				Usage:   "Skip TLS certificate verification (development only)",
				Sources: cli.EnvVars("PARSE_DMARC_MCP_OAUTH_INSECURE"),
			},
			&cli.BoolFlag{
				Name:    "oauth-login",
				Usage:   "Run the IMAP OAuth2 device flow, save the refresh token, and exit",
				Sources: cli.EnvVars("PARSE_DMARC_OAUTH_LOGIN"),
			},
		},
		Action: run,
		Commands: []*cli.Command{
			{
				Name:  "version",
				Usage: "Show version information",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					fmt.Printf("Version:    %s\n", version)
					fmt.Printf("Commit:     %s\n", commit)
					fmt.Printf("Build Date: %s\n", date)
					fmt.Printf("Built By:   %s\n", builtBy)
					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal().Err(err).Msg("failed to run")
	}
}

func run(ctx context.Context, cmd *cli.Command) error {
	configPath := cmd.String("config")
	genConfig := cmd.Bool("gen-config")
	fetchOnce := cmd.Bool("fetch-once")
	serveOnly := cmd.Bool("serve-only")
	fetchInterval := cmd.Int("fetch-interval")
	metricsEnabled := cmd.Bool("metrics")
	mcpMode := cmd.Bool("mcp")
	mcpHTTPAddr := cmd.String("mcp-http")

	// OAuth configuration for MCP HTTP server
	mcpOAuthEnabled := cmd.Bool("mcp-oauth")
	mcpOAuthIssuer := cmd.String("mcp-oauth-issuer")
	mcpOAuthAudience := cmd.String("mcp-oauth-audience")
	mcpOAuthClientID := cmd.String("mcp-oauth-client-id")
	mcpOAuthClientSecret := cmd.String("mcp-oauth-client-secret")
	mcpOAuthScopes := cmd.String("mcp-oauth-scopes")
	mcpOAuthIntrospection := cmd.String("mcp-oauth-introspection-endpoint")
	mcpOAuthResourceName := cmd.String("mcp-oauth-resource-name")
	mcpOAuthInsecure := cmd.Bool("mcp-oauth-insecure")

	if genConfig {
		if err := config.GenerateSample(configPath); err != nil {
			return fmt.Errorf("failed to generate config: %w", err)
		}
		log.Info().Str("path", configPath).Msg("sample configuration written")
		return nil
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Reinitialize logger with config-derived level
	log = logger.NewLogger(cfg.LogLevel, !cfg.ColoredLogs)

	if cmd.Bool("oauth-login") {
		return runOAuthLogin(ctx, cfg)
	}

	// Validate required IMAP configuration when fetching is enabled
	// (not serve-only and not MCP mode)
	if !serveOnly && !mcpMode && mcpHTTPAddr == "" {
		if err := cfg.Validate(); err != nil {
			return fmt.Errorf("configuration error: %w", err)
		}
	}

	if err := os.MkdirAll(filepath.Dir(cfg.Database.Path), 0o755); err != nil {
		return fmt.Errorf("create database directory: %w", err)
	}

	store, err := storage.NewStorage(cfg.Database.Path)
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}
	defer func() { _ = store.Close() }()

	// Handle MCP mode
	if mcpMode || mcpHTTPAddr != "" {
		// Build OAuth config if enabled
		var oauthCfg *oauth.Config
		if mcpOAuthEnabled {
			// Parse scopes
			var scopes []string
			if mcpOAuthScopes != "" {
				for _, s := range strings.Split(mcpOAuthScopes, ",") {
					scopes = append(scopes, strings.TrimSpace(s))
				}
			}

			// Determine resource server URL and audience from HTTP address
			var resourceServerURL, audience string
			if mcpOAuthAudience != "" {
				resourceServerURL = mcpOAuthAudience
				audience = mcpOAuthAudience
			} else if mcpHTTPAddr != "" {
				// Use localhost with the port if no audience specified
				resourceServerURL = "http://localhost" + mcpHTTPAddr
				audience = resourceServerURL
			}

			oauthCfg = &oauth.Config{
				Enabled:               true,
				Issuer:                mcpOAuthIssuer,
				Audience:              audience,
				ClientID:              mcpOAuthClientID,
				ClientSecret:          mcpOAuthClientSecret,
				RequiredScopes:        scopes,
				IntrospectionEndpoint: mcpOAuthIntrospection,
				ResourceServerURL:     resourceServerURL,
				ResourceName:          mcpOAuthResourceName,
				InsecureSkipVerify:    mcpOAuthInsecure,
			}
		}
		return runMCPServer(ctx, store, mcpHTTPAddr, oauthCfg)
	}

	// Initialize metrics if enabled
	var m *metrics.Metrics
	if metricsEnabled {
		m = metrics.New(version, commit, date)
		log.Info().Msg("prometheus metrics enabled at /metrics")
	}

	// Build IMAP token source if XOAUTH2 is configured. Done once so the
	// underlying ReuseTokenSource caches access tokens across fetch cycles.
	var tokenSource oauth2.TokenSource
	if cfg.IMAP.Auth.IsXOAUTH2() && !serveOnly {
		tokenSource, err = buildTokenSource(ctx, cfg)
		if err != nil {
			return fmt.Errorf("initialize oauth token source: %w", err)
		}
	}

	ctx, stop := signal.NotifyContext(ctx, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	server := api.NewServer(store, cfg.Server.Host, cfg.Server.Port, m, log)
	serverErrChan := make(chan error, 1)
	go func() {
		serverErrChan <- server.Start(ctx)
	}()

	// Refresh metrics on startup
	server.RefreshMetrics()

	if serveOnly {
		log.Info().Msg("running in serve-only mode")
		select {
		case <-ctx.Done():
			log.Info().Msg("shutting down")
		case err := <-serverErrChan:
			if err != nil {
				return fmt.Errorf("server error: %w", err)
			}
		}
		return nil
	}

	if fetchOnce {
		if err := fetchReports(cfg, store, m, tokenSource); err != nil {
			return fmt.Errorf("failed to fetch reports: %w", err)
		}
		server.RefreshMetrics()
		log.Info().Msg("fetch complete")
		return nil
	}

	log.Info().Int("interval_seconds", fetchInterval).Msg("starting continuous fetch mode")

	if err := fetchReports(cfg, store, m, tokenSource); err != nil {
		handleFetchError(err, m)
	}
	server.RefreshMetrics()

	ticker := time.NewTicker(time.Duration(fetchInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := fetchReports(cfg, store, m, tokenSource); err != nil {
				handleFetchError(err, m)
			}
			server.RefreshMetrics()
		case <-ctx.Done():
			log.Info().Msg("shutting down")
			return nil
		case err := <-serverErrChan:
			if err != nil {
				return fmt.Errorf("server error: %w", err)
			}
		}
	}
}

func fetchReports(cfg *config.Config, store *storage.Storage, m *metrics.Metrics, tokenSource oauth2.TokenSource) error {
	log.Info().Msg("fetching DMARC reports")

	fetchStart := time.Now()
	if m != nil {
		m.FetchCyclesTotal.Inc()
	}

	// Create IMAP client
	connectStart := time.Now()
	client := imap.NewClient(&cfg.IMAP, log, tokenSource)
	if err := client.Connect(); err != nil {
		if m != nil {
			m.RecordIMAPConnection(false, time.Since(connectStart))
			m.FetchErrors.Inc()
		}
		return fmt.Errorf("connect to IMAP server: %w", err)
	}
	if m != nil {
		m.RecordIMAPConnection(true, time.Since(connectStart))
		m.IMAPAuthRequired.Set(0)
	}
	defer func() { _ = client.Disconnect() }()

	// Fetch reports
	result, err := client.FetchDMARCReports()
	if err != nil {
		if m != nil {
			m.FetchErrors.Inc()
		}
		return fmt.Errorf("fetch DMARC reports: %w", err)
	}

	if m != nil {
		m.ReportsFetched.Add(float64(len(result.Reports)))
	}

	if len(result.Reports) == 0 {
		log.Info().Msg("no new reports found")
		if m != nil {
			m.RecordFetchDuration(time.Since(fetchStart))
			m.LastFetchTimestamp.SetToCurrentTime()
		}
		return nil
	}

	log.Info().Int("count", len(result.Reports)).Msg("processing reports")

	// Process each report
	processed := 0
	for _, report := range result.Reports {
		for _, attachment := range report.Attachments {
			if m != nil {
				m.AttachmentsTotal.Inc()
			}

			feedback, err := parser.ParseReport(attachment.Data)
			if err != nil {
				log.Warn().Err(err).Str("filename", attachment.Filename).Msg("failed to parse report")
				if m != nil {
					m.ReportParseErrors.Inc()
				}
				continue
			}
			if m != nil {
				m.ReportsParsed.Inc()
			}

			if err := store.SaveReport(feedback); err != nil {
				log.Error().Err(err).Str("report_id", feedback.ReportMetadata.ReportID).Msg("failed to save report")
				if m != nil {
					m.ReportStoreErrors.Inc()
				}
				continue
			}
			if m != nil {
				m.ReportsStored.Inc()
			}

			log.Info().
				Str("report_id", feedback.ReportMetadata.ReportID).
				Str("org", feedback.ReportMetadata.OrgName).
				Str("domain", feedback.PolicyPublished.Domain).
				Int("messages", feedback.GetTotalMessages()).
				Msg("saved report")
			processed++
		}
	}

	// Post-processing: mark as seen and/or move messages
	if len(result.MessageIDs) > 0 {
		if cfg.IMAP.MarkAsSeen {
			if err := client.MarkAsSeen(result.MessageIDs); err != nil {
				log.Error().Err(err).Msg("failed to mark messages as seen")
			} else {
				log.Info().Int("count", len(result.MessageIDs)).Msg("marked messages as seen")
			}
		}
		if cfg.IMAP.ProcessedMailbox != "" {
			if err := client.MoveMessages(result.MessageIDs, cfg.IMAP.ProcessedMailbox); err != nil {
				log.Error().Err(err).Msg("failed to move messages to processed mailbox")
			} else {
				log.Info().Int("count", len(result.MessageIDs)).Str("mailbox", cfg.IMAP.ProcessedMailbox).Msg("moved messages to processed mailbox")
			}
		}
	}

	if m != nil {
		m.RecordFetchDuration(time.Since(fetchStart))
		m.LastFetchTimestamp.SetToCurrentTime()
	}

	log.Info().Int("count", processed).Msg("reports processed")
	return nil
}

// buildTokenSource constructs an XOAUTH2 token source backed by the refresh
// token in secrets.json (or the IMAP_OAUTH_REFRESH_TOKEN env override).
func buildTokenSource(ctx context.Context, cfg *config.Config) (oauth2.TokenSource, error) {
	provider, err := imapoauth.ProviderByName(cfg.IMAP.Auth.Provider)
	if err != nil {
		return nil, err
	}
	refreshToken, envOverride, err := imapoauth.LoadSecrets(cfg.SecretsPath())
	if err != nil {
		return nil, err
	}
	if envOverride {
		log.Info().Msg("using oauth refresh token from IMAP_OAUTH_REFRESH_TOKEN env (read-only)")
	} else {
		log.Info().Str("path", cfg.SecretsPath()).Msg("loaded oauth refresh token from secrets file")
	}
	return imapoauth.NewTokenSource(ctx, provider, cfg.IMAP.Auth.ClientID, cfg.IMAP.Auth.ClientSecret, refreshToken), nil
}

// runOAuthLogin runs the OAuth2 device authorization grant and writes the
// resulting refresh token to secrets.json. Exits the process on completion.
func runOAuthLogin(ctx context.Context, cfg *config.Config) error {
	if !cfg.IMAP.Auth.IsXOAUTH2() {
		return fmt.Errorf("oauth-login requires imap.auth.type=xoauth2 in config (or IMAP_AUTH_TYPE=xoauth2)")
	}
	if cfg.IMAP.Auth.ClientID == "" || cfg.IMAP.Auth.ClientSecret == "" {
		return fmt.Errorf("imap.auth.client_id and client_secret must be set before --oauth-login")
	}

	provider, err := imapoauth.ProviderByName(cfg.IMAP.Auth.Provider)
	if err != nil {
		return err
	}

	prompt := func(authURL, _ string) {
		fmt.Println()
		fmt.Println("Open this URL in a browser to authorize dmarcguard:")
		fmt.Println("  " + authURL)
		fmt.Println()
		fmt.Println("(attempting to open it for you automatically)")
		fmt.Println("Waiting for callback on http://127.0.0.1:<port>/callback ...")
	}

	result, err := imapoauth.LoopbackLogin(ctx, provider, cfg.IMAP.Auth.ClientID, cfg.IMAP.Auth.ClientSecret, prompt)
	if err != nil {
		return fmt.Errorf("loopback login: %w", err)
	}

	fmt.Println()
	if result.Email != "" {
		fmt.Println("✓ Authorized as " + result.Email)
		if cfg.IMAP.Username != "" && cfg.IMAP.Username != result.Email {
			fmt.Println()
			fmt.Println("⚠ Your config sets imap.username=" + cfg.IMAP.Username)
			fmt.Println("  but you authenticated as " + result.Email)
			fmt.Println("  Gmail will reject the XOAUTH2 connection unless these match.")
		}
	}

	// Try to write the secrets file at the configured path. This works when
	// running on the same host as the daemon. When the path isn't writable
	// (e.g. the binary is run on a developer Mac but config.json points at a
	// Docker path like /data/secrets.json), fall back to printing the token
	// so the user can paste it into IMAP_OAUTH_REFRESH_TOKEN.
	if err := imapoauth.SaveSecrets(cfg.SecretsPath(), result.RefreshToken); err != nil {
		fmt.Println()
		fmt.Println("⚠ Could not write " + cfg.SecretsPath() + ": " + err.Error())
		fmt.Println("  This is expected if you're running --oauth-login on a host but the daemon")
		fmt.Println("  runs in Docker (the configured secrets path doesn't exist on the host).")
		fmt.Println()
		fmt.Println("Add this env var to your compose.yml or run command:")
		fmt.Println()
		fmt.Println("    IMAP_OAUTH_REFRESH_TOKEN=" + result.RefreshToken)
		fmt.Println()
		return nil
	}

	fmt.Println("✓ Refresh token saved to " + cfg.SecretsPath())
	return nil
}

// handleFetchError logs a fetch error and, when it represents a permanent
// OAuth failure, sets the auth-required gauge so operators see a clear signal.
func handleFetchError(err error, m *metrics.Metrics) {
	if imapoauth.IsTerminalAuthError(err) {
		log.Error().Err(err).Msg("oauth refresh token rejected — re-run with --oauth-login")
		if m != nil {
			m.IMAPAuthRequired.Set(1)
		}
		return
	}
	log.Error().Err(err).Msg("fetch failed")
}

func runMCPServer(ctx context.Context, store *storage.Storage, httpAddr string, oauthCfg *oauth.Config) error {
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	mcpCfg := &mcpserver.Config{
		Version:  version,
		HTTPAddr: httpAddr,
		Logger:   log,
		OAuth:    oauthCfg,
	}

	server := mcpserver.NewServer(store, mcpCfg)

	// If HTTP address is specified, run HTTP server
	// Otherwise, run over stdio
	if httpAddr != "" {
		return server.RunHTTP(ctx, mcpCfg.HTTPAddr, oauthCfg)
	}

	return server.RunStdio(ctx)
}
