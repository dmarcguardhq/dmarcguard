package api

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/goccy/go-json"

	"github.com/meysam81/parse-dmarc/internal/metrics"
	"github.com/meysam81/parse-dmarc/internal/storage"
)

//go:embed dist
var distFS embed.FS

// Server represents the API server
type Server struct {
	storage *storage.Storage
	metrics *metrics.Metrics
	addr    string
}

// NewServer creates a new API server
func NewServer(store *storage.Storage, host string, port int, m *metrics.Metrics) *Server {
	return &Server{
		storage: store,
		metrics: m,
		addr:    fmt.Sprintf("%s:%d", host, port),
	}
}

// Start starts the HTTP server
func (s *Server) Start(ctx context.Context) error {
	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("/api/reports", s.handleReports)
	mux.HandleFunc("/api/reports/", s.handleReportDetail)
	mux.HandleFunc("/api/statistics", s.handleStatistics)
	mux.HandleFunc("/api/top-sources", s.handleTopSources)

	// Prometheus metrics endpoint
	if s.metrics != nil {
		mux.Handle("/metrics", s.metrics.Handler())
	}

	// Serve frontend
	// Try to serve embedded files, fallback to nothing if not embedded
	distFiles, err := fs.Sub(distFS, "dist")
	if err == nil {
		mux.Handle("/", http.FileServer(http.FS(distFiles)))
	} else {
		// If dist folder is not embedded, serve a simple message
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
		log.Println("Shutting down server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("Server shutdown error: %v", err)
		}
	}()

	log.Printf("Starting server on %s", s.addr)
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// corsMiddleware adds CORS headers
func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

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

	// Parse pagination parameters
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

	// Extract ID from URL
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

// writeJSON writes JSON response
func (s *Server) writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Failed to encode JSON: %v", err)
	}
}

// RefreshMetrics updates all Prometheus metrics from current database state
func (s *Server) RefreshMetrics() {
	if s.metrics == nil {
		return
	}

	// Update basic statistics
	stats, err := s.storage.GetStatistics()
	if err != nil {
		log.Printf("Failed to get statistics for metrics: %v", err)
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

	// Update per-domain metrics
	domainStats, err := s.storage.GetDomainStats()
	if err != nil {
		log.Printf("Failed to get domain stats for metrics: %v", err)
	} else {
		for _, ds := range domainStats {
			s.metrics.UpdateDomainMetrics(ds.Domain, ds.TotalMessages, ds.ComplianceRate)
		}
	}

	// Update per-organization metrics
	orgStats, err := s.storage.GetOrgStats()
	if err != nil {
		log.Printf("Failed to get org stats for metrics: %v", err)
	} else {
		for _, os := range orgStats {
			s.metrics.UpdateOrgMetrics(os.OrgName, os.Reports)
		}
	}

	// Update disposition metrics
	dispStats, err := s.storage.GetDispositionStats()
	if err != nil {
		log.Printf("Failed to get disposition stats for metrics: %v", err)
	} else {
		for _, ds := range dispStats {
			s.metrics.UpdateDispositionMetrics(ds.Disposition, ds.Count)
		}
	}

	// Update authentication results
	spfStats, errSpf := s.storage.GetSPFStats()
	dkimStats, errDkim := s.storage.GetDKIMStats()
	if errSpf != nil {
		log.Printf("Failed to get SPF stats for metrics: %v", errSpf)
	}
	if errDkim != nil {
		log.Printf("Failed to get DKIM stats for metrics: %v", errDkim)
	}
	if errSpf == nil && errDkim == nil {
		spfResults := make(map[string]int)
		for _, s := range spfStats {
			spfResults[s.Result] = s.Count
		}
		dkimResults := make(map[string]int)
		for _, d := range dkimStats {
			dkimResults[d.Result] = d.Count
		}
		s.metrics.UpdateAuthResults(spfResults, dkimResults)
	}
}

// GetMetrics returns the metrics instance
func (s *Server) GetMetrics() *metrics.Metrics {
	return s.metrics
}
