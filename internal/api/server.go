package api

import (
	"context"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/goccy/go-json"
	"github.com/rs/zerolog"

	mcpserver "github.com/meysam81/parse-dmarc/internal/mcp"
	"github.com/meysam81/parse-dmarc/internal/mcp/oauth"
	"github.com/meysam81/parse-dmarc/internal/metrics"
	"github.com/meysam81/parse-dmarc/internal/storage"
)

//go:embed dist
var distFS embed.FS

const (
	// SettingMCPEnabled is the database key for MCP enabled state.
	SettingMCPEnabled = "mcp_enabled"
)

// MCPConfig holds MCP server configuration.
type MCPConfig struct {
	Version string
	OAuth   *oauth.Config
	Logger  *zerolog.Logger
}

// Server represents the API server
type Server struct {
	storage *storage.Storage
	metrics *metrics.Metrics
	addr    string
	logger  *zerolog.Logger

	// MCP integration
	mcpConfig  *MCPConfig
	mcpServer  *mcpserver.Server
	mcpHandler http.Handler
	mcpMu      sync.RWMutex
	mcpEnabled bool
}

// NewServer creates a new API server
func NewServer(store *storage.Storage, host string, port int, m *metrics.Metrics) *Server {
	return &Server{
		storage: store,
		metrics: m,
		addr:    fmt.Sprintf("%s:%d", host, port),
	}
}

// SetLogger sets the logger for the server.
func (s *Server) SetLogger(logger *zerolog.Logger) {
	s.logger = logger
}

// SetMCPConfig configures the MCP server integration.
func (s *Server) SetMCPConfig(cfg *MCPConfig) {
	s.mcpConfig = cfg
}

// initMCP initializes the MCP server and handler.
func (s *Server) initMCP() error {
	if s.mcpConfig == nil {
		return nil
	}

	s.mcpMu.Lock()
	defer s.mcpMu.Unlock()

	// Check if MCP is enabled in settings (default: disabled)
	enabled, err := s.storage.GetSetting(SettingMCPEnabled)
	if err != nil {
		return fmt.Errorf("failed to get MCP enabled setting: %w", err)
	}
	s.mcpEnabled = enabled == "true"

	// Create MCP server
	mcpCfg := &mcpserver.Config{
		Version: s.mcpConfig.Version,
		Logger:  s.mcpConfig.Logger,
		OAuth:   s.mcpConfig.OAuth,
	}
	s.mcpServer = mcpserver.NewServer(s.storage, mcpCfg)

	// Create handler
	handler, err := s.mcpServer.Handler("/mcp", s.mcpConfig.OAuth)
	if err != nil {
		return fmt.Errorf("failed to create MCP handler: %w", err)
	}
	s.mcpHandler = handler

	return nil
}

// IsMCPEnabled returns the current MCP enabled state.
func (s *Server) IsMCPEnabled() bool {
	s.mcpMu.RLock()
	defer s.mcpMu.RUnlock()
	return s.mcpEnabled
}

// SetMCPEnabled enables or disables the MCP server.
func (s *Server) SetMCPEnabled(enabled bool) error {
	s.mcpMu.Lock()
	defer s.mcpMu.Unlock()

	value := "false"
	if enabled {
		value = "true"
	}

	if err := s.storage.SetSetting(SettingMCPEnabled, value); err != nil {
		return fmt.Errorf("failed to save MCP enabled setting: %w", err)
	}

	s.mcpEnabled = enabled

	if s.logger != nil {
		s.logger.Info().Bool("enabled", enabled).Msg("MCP server state changed")
	}

	return nil
}

// Start starts the HTTP server
func (s *Server) Start(ctx context.Context) error {
	// Initialize MCP if configured
	if err := s.initMCP(); err != nil {
		if s.logger != nil {
			s.logger.Error().Err(err).Msg("failed to initialize MCP server")
		}
	}

	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("/api/reports", s.handleReports)
	mux.HandleFunc("/api/reports/", s.handleReportDetail)
	mux.HandleFunc("/api/statistics", s.handleStatistics)
	mux.HandleFunc("/api/top-sources", s.handleTopSources)
	mux.HandleFunc("/api/settings", s.handleSettings)
	mux.HandleFunc("/api/settings/mcp", s.handleMCPSetting)

	// Prometheus metrics endpoint
	if s.metrics != nil {
		mux.Handle("/metrics", s.metrics.Handler())
	}

	// MCP endpoint - dynamically enabled/disabled
	if s.mcpHandler != nil {
		mux.HandleFunc("/mcp/", s.handleMCP)
		mux.HandleFunc("/mcp", s.handleMCP)
	}

	// Serve frontend
	distFiles, err := fs.Sub(distFS, "dist")
	if err == nil {
		fileServer := http.FileServer(http.FS(distFiles))
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			// Don't serve static files for API or MCP routes
			if strings.HasPrefix(r.URL.Path, "/api/") || strings.HasPrefix(r.URL.Path, "/mcp") {
				http.NotFound(w, r)
				return
			}
			fileServer.ServeHTTP(w, r)
		})
	} else {
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/" {
				w.Header().Set("Content-Type", "text/html")
				_, _ = fmt.Fprintf(w, `
					<!DOCTYPE html>
					<html>
					<head><title>DMARC Dashboard</title></head>
					<body>
						<h1>DMARC Report Dashboard API</h1>
						<p>API is running. Frontend assets not embedded yet.</p>
						<ul>
							<li><a href="/api/statistics">Statistics</a></li>
							<li><a href="/api/reports">Reports</a></li>
							<li><a href="/api/top-sources">Top Sources</a></li>
							<li><a href="/metrics">Prometheus Metrics</a></li>
						</ul>
					</body>
					</html>
				`)
			} else {
				http.NotFound(w, r)
			}
		})
	}

	// Build handler chain: CORS -> Metrics -> Routes
	var handler http.Handler = mux
	if s.metrics != nil {
		handler = s.metrics.HTTPMiddleware(handler)
	}
	handler = s.corsMiddleware(handler)

	server := &http.Server{
		Addr:    s.addr,
		Handler: handler,
	}

	go func() {
		<-ctx.Done()
		if s.logger != nil {
			s.logger.Info().Msg("shutting down server")
		}
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			if s.logger != nil {
				s.logger.Error().Err(err).Msg("server shutdown error")
			}
		}
	}()

	if s.logger != nil {
		s.logger.Info().Str("addr", s.addr).Msg("starting server")
	}
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// handleMCP proxies requests to the MCP handler when enabled.
func (s *Server) handleMCP(w http.ResponseWriter, r *http.Request) {
	if !s.IsMCPEnabled() {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte(`{"error":"MCP server is disabled","message":"Enable MCP in settings to use this endpoint"}`))
		return
	}

	s.mcpMu.RLock()
	handler := s.mcpHandler
	s.mcpMu.RUnlock()

	if handler == nil {
		http.Error(w, "MCP server not configured", http.StatusInternalServerError)
		return
	}

	// Strip /mcp prefix for the MCP handler
	r.URL.Path = strings.TrimPrefix(r.URL.Path, "/mcp")
	if r.URL.Path == "" {
		r.URL.Path = "/"
	}

	handler.ServeHTTP(w, r)
}

// corsMiddleware adds CORS headers
func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// handleReports returns a list of reports
func (s *Server) handleReports(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	limit := 50
	offset := 0

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	reports, err := s.storage.GetReports(limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.writeJSON(w, reports)
}

// handleReportDetail returns a single report detail
func (s *Server) handleReportDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Path[len("/api/reports/"):]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid report ID", http.StatusBadRequest)
		return
	}

	report, err := s.storage.GetReportByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	s.writeJSON(w, report)
}

// handleStatistics returns dashboard statistics
func (s *Server) handleStatistics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	stats, err := s.storage.GetStatistics()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.writeJSON(w, stats)
}

// handleTopSources returns top source IPs
func (s *Server) handleTopSources(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	limit := 10
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	sources, err := s.storage.GetTopSourceIPs(limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.writeJSON(w, sources)
}

// SettingsResponse represents the settings API response.
type SettingsResponse struct {
	MCPEnabled bool   `json:"mcp_enabled"`
	MCPPath    string `json:"mcp_path"`
}

// handleSettings returns all settings
func (s *Server) handleSettings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := SettingsResponse{
		MCPEnabled: s.IsMCPEnabled(),
		MCPPath:    "/mcp",
	}

	s.writeJSON(w, response)
}

// MCPSettingRequest represents a request to update MCP setting.
type MCPSettingRequest struct {
	Enabled bool `json:"enabled"`
}

// MCPSettingResponse represents the MCP setting response.
type MCPSettingResponse struct {
	Enabled bool   `json:"enabled"`
	Path    string `json:"path"`
	Message string `json:"message,omitempty"`
}

// handleMCPSetting handles GET and PUT for MCP enabled setting
func (s *Server) handleMCPSetting(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		response := MCPSettingResponse{
			Enabled: s.IsMCPEnabled(),
			Path:    "/mcp",
		}
		s.writeJSON(w, response)

	case http.MethodPut:
		if s.mcpHandler == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error":"MCP server not configured"}`))
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var req MCPSettingRequest
		if err := json.Unmarshal(body, &req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if err := s.SetMCPEnabled(req.Enabled); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		msg := "MCP server disabled"
		if req.Enabled {
			msg = "MCP server enabled"
		}

		response := MCPSettingResponse{
			Enabled: req.Enabled,
			Path:    "/mcp",
			Message: msg,
		}
		s.writeJSON(w, response)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// writeJSON writes JSON response
func (s *Server) writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		if s.logger != nil {
			s.logger.Error().Err(err).Msg("failed to encode JSON")
		}
	}
}

// RefreshMetrics updates all Prometheus metrics from current database state
func (s *Server) RefreshMetrics() {
	if s.metrics == nil {
		return
	}

	stats, err := s.storage.GetStatistics()
	if err != nil {
		if s.logger != nil {
			s.logger.Error().Err(err).Msg("failed to get statistics for metrics")
		}
	} else {
		s.metrics.UpdateStatistics(
			stats.TotalReports,
			stats.TotalMessages,
			stats.CompliantMessages,
			stats.UniqueSourceIPs,
			stats.UniqueDomains,
			stats.ComplianceRate,
		)
	}

	domainStats, err := s.storage.GetDomainStats()
	if err != nil {
		if s.logger != nil {
			s.logger.Error().Err(err).Msg("failed to get domain stats for metrics")
		}
	} else {
		for _, ds := range domainStats {
			s.metrics.UpdateDomainMetrics(ds.Domain, ds.TotalMessages, ds.ComplianceRate)
		}
	}

	orgStats, err := s.storage.GetOrgStats()
	if err != nil {
		if s.logger != nil {
			s.logger.Error().Err(err).Msg("failed to get org stats for metrics")
		}
	} else {
		for _, os := range orgStats {
			s.metrics.UpdateOrgMetrics(os.OrgName, os.Reports)
		}
	}

	dispStats, err := s.storage.GetDispositionStats()
	if err != nil {
		if s.logger != nil {
			s.logger.Error().Err(err).Msg("failed to get disposition stats for metrics")
		}
	} else {
		for _, ds := range dispStats {
			s.metrics.UpdateDispositionMetrics(ds.Disposition, ds.Count)
		}
	}

	spfStats, errSpf := s.storage.GetSPFStats()
	dkimStats, errDkim := s.storage.GetDKIMStats()
	if errSpf != nil {
		if s.logger != nil {
			s.logger.Error().Err(errSpf).Msg("failed to get SPF stats for metrics")
		}
	}
	if errDkim != nil {
		if s.logger != nil {
			s.logger.Error().Err(errDkim).Msg("failed to get DKIM stats for metrics")
		}
	}
	if errSpf == nil && errDkim == nil {
		spfResults := make(map[string]int)
		for _, stat := range spfStats {
			spfResults[stat.Result] = stat.Count
		}
		dkimResults := make(map[string]int)
		for _, stat := range dkimStats {
			dkimResults[stat.Result] = stat.Count
		}
		s.metrics.UpdateAuthResults(spfResults, dkimResults)
	}
}

// GetMetrics returns the metrics instance
func (s *Server) GetMetrics() *metrics.Metrics {
	return s.metrics
}
