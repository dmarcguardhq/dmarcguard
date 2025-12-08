package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/meysam81/parse-dmarc/internal/api"
	"github.com/meysam81/parse-dmarc/internal/config"
	"github.com/meysam81/parse-dmarc/internal/imap"
	"github.com/meysam81/parse-dmarc/internal/logger"
	mcpserver "github.com/meysam81/parse-dmarc/internal/mcp"
	"github.com/meysam81/parse-dmarc/internal/metrics"
	"github.com/meysam81/parse-dmarc/internal/parser"
	"github.com/meysam81/parse-dmarc/internal/storage"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v3"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"

	log         *zerolog.Logger
	logLevel    = "info"
	coloredLogs = false
)

func main() {
	cmd := &cli.Command{
		Name:                  "parse-dmarc",
		Usage:                 "Parse and analyze DMARC reports from IMAP mailbox",
		Version:               version,
		EnableShellCompletion: true,
		Suggest:               true,
		Before: func(ctx context.Context, c *cli.Command) (context.Context, error) {
			log = logger.NewLogger(logLevel, !coloredLogs)
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

	if genConfig {
		if err := config.GenerateSample(configPath); err != nil {
			return fmt.Errorf("failed to generate config: %w", err)
		}
		log.Printf("Sample configuration written to %s", configPath)
		return nil
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	store, err := storage.NewStorage(cfg.Database.Path)
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}
	defer func() { _ = store.Close() }()

	// Handle MCP mode
	if mcpMode || mcpHTTPAddr != "" {
		return runMCPServer(ctx, store, mcpHTTPAddr)
	}

	// Initialize metrics if enabled
	var m *metrics.Metrics
	if metricsEnabled {
		m = metrics.New(version, commit, date)
		log.Println("Prometheus metrics enabled at /metrics")
	}

	ctx, stop := signal.NotifyContext(ctx, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	server := api.NewServer(store, cfg.Server.Host, cfg.Server.Port, m)
	serverErrChan := make(chan error, 1)
	go func() {
		serverErrChan <- server.Start(ctx)
	}()

	// Refresh metrics on startup
	server.RefreshMetrics()

	if serveOnly {
		log.Println("Running in serve-only mode")
		select {
		case <-ctx.Done():
			log.Println("Shutting down...")
		case err := <-serverErrChan:
			if err != nil {
				return fmt.Errorf("server error: %w", err)
			}
		}
		return nil
	}

	if fetchOnce {
		if err := fetchReports(cfg, store, m); err != nil {
			return fmt.Errorf("failed to fetch reports: %w", err)
		}
		server.RefreshMetrics()
		log.Println("Fetch complete")
		return nil
	}

	log.Printf("Starting continuous fetch mode (interval: %d seconds)", fetchInterval)

	if err := fetchReports(cfg, store, m); err != nil {
		log.Printf("Initial fetch failed: %v", err)
	}
	server.RefreshMetrics()

	ticker := time.NewTicker(time.Duration(fetchInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := fetchReports(cfg, store, m); err != nil {
				log.Printf("Fetch failed: %v", err)
			}
			server.RefreshMetrics()
		case <-ctx.Done():
			log.Println("Shutting down...")
			return nil
		case err := <-serverErrChan:
			if err != nil {
				return fmt.Errorf("server error: %w", err)
			}
		}
	}
}

func fetchReports(cfg *config.Config, store *storage.Storage, m *metrics.Metrics) error {
	log.Println("Fetching DMARC reports...")

	fetchStart := time.Now()
	if m != nil {
		m.FetchCyclesTotal.Inc()
	}

	// Create IMAP client
	connectStart := time.Now()
	client := imap.NewClient(&cfg.IMAP)
	if err := client.Connect(); err != nil {
		if m != nil {
			m.RecordIMAPConnection(false, time.Since(connectStart))
			m.FetchErrors.Inc()
		}
		return err
	}
	if m != nil {
		m.RecordIMAPConnection(true, time.Since(connectStart))
	}
	defer func() { _ = client.Disconnect() }()

	// Fetch reports
	reports, err := client.FetchDMARCReports()
	if err != nil {
		if m != nil {
			m.FetchErrors.Inc()
		}
		return err
	}

	if m != nil {
		m.ReportsFetched.Add(float64(len(reports)))
	}

	if len(reports) == 0 {
		log.Println("No new reports found")
		if m != nil {
			m.RecordFetchDuration(time.Since(fetchStart))
			m.LastFetchTimestamp.SetToCurrentTime()
		}
		return nil
	}

	log.Printf("Processing %d reports...", len(reports))

	// Process each report
	processed := 0
	for _, report := range reports {
		for _, attachment := range report.Attachments {
			if m != nil {
				m.AttachmentsTotal.Inc()
			}

			feedback, err := parser.ParseReport(attachment.Data)
			if err != nil {
				log.Printf("Failed to parse %s: %v", attachment.Filename, err)
				if m != nil {
					m.ReportParseErrors.Inc()
				}
				continue
			}
			if m != nil {
				m.ReportsParsed.Inc()
			}

			if err := store.SaveReport(feedback); err != nil {
				log.Printf("Failed to save report %s: %v", feedback.ReportMetadata.ReportID, err)
				if m != nil {
					m.ReportStoreErrors.Inc()
				}
				continue
			}
			if m != nil {
				m.ReportsStored.Inc()
			}

			log.Printf("Saved report: %s from %s (domain: %s, messages: %d)",
				feedback.ReportMetadata.ReportID,
				feedback.ReportMetadata.OrgName,
				feedback.PolicyPublished.Domain,
				feedback.GetTotalMessages())
			processed++
		}
	}

	if m != nil {
		m.RecordFetchDuration(time.Since(fetchStart))
		m.LastFetchTimestamp.SetToCurrentTime()
	}

	log.Printf("Successfully processed %d reports", processed)
	return nil
}

func runMCPServer(ctx context.Context, store *storage.Storage, httpAddr string) error {
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	mcpCfg := &mcpserver.Config{
		Version:  version,
		HTTPAddr: httpAddr,
		Logger:   log,
	}

	server := mcpserver.NewServer(store, mcpCfg)

	// If HTTP address is specified, run HTTP server
	// Otherwise, run over stdio
	if httpAddr != "" {
		return server.RunHTTP(ctx, mcpCfg.HTTPAddr)
	}

	return server.RunStdio(ctx)
}
