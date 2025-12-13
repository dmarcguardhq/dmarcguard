# Parse DMARC - Project Guide

A Go application that fetches DMARC reports from IMAP mailboxes, parses them, and displays them in a Vue.js dashboard.

## Tech Stack

- **Backend**: Go 1.24.7
- **Frontend**: Vue.js 3 with Vite
- **Database**: SQLite (supports both CGO and pure-Go variants)
- **Package Manager**: Bun (for frontend)
- **Task Runner**: Just (Justfile)
- **CLI Framework**: urfave/cli/v3
- **JSON Library**: goccy/go-json (high-performance)

## Project Structure

```
parse-dmarc/
├── main.go                # Main application entry point
├── internal/
│   ├── api/               # REST API server and embedded frontend
│   ├── config/            # Configuration management (JSON + env vars)
│   ├── imap/              # IMAP client for fetching emails
│   ├── parser/            # DMARC XML parser
│   └── storage/           # SQLite database layer (CGO + pure-Go)
├── src/                   # Vue.js 3 frontend source
│   ├── App.vue            # Main application component (dashboard)
│   └── components/
│       └── DnsGenerator.vue  # DMARC DNS record generator tool
├── public/                # Static frontend assets
├── assets/                # Project assets (icons, images)
├── scripts/               # Utility scripts (SVG to PNG conversion)
├── Justfile               # Build commands
├── Dockerfile             # Multi-stage Docker build
├── ROADMAP.md             # Product roadmap and feature planning
├── CONTRIBUTING.md        # Contribution guidelines
└── .goreleaser.yaml       # Release automation
```

## Development Commands

```bash
# Install all dependencies (Go + Node)
just install-deps

# Build full application (frontend + backend)
just build

# Build with CGO (for native SQLite)
just build-cgo

# Run development server with hot reload (uses air)
just dev

# Run frontend dev server only
just frontend-dev

# Run tests
just test

# Generate config file
just config

# Clean build artifacts
just clean

# Install binary to /usr/local/bin
just install
```

## Building

The build process:

1. Frontend is built with `bun run build`
2. Frontend dist is copied to `internal/api/dist/`
3. Go binary embeds the frontend and serves it
4. Output binary: `bin/parse-dmarc`

Two build modes:

- `just build` - Pure Go (CGO_ENABLED=0), uses modernc.org/sqlite
- `just build-cgo` - CGO enabled, uses mattn/go-sqlite3

## Testing

```bash
# Run all Go tests
go test -v ./...

# Run tests for specific package
go test -v ./internal/parser/...
```

## Code Style

- Go: Standard gofmt formatting, golangci-lint for linting
- Frontend: Vue 3 Composition API
- Pre-commit hooks configured in `.pre-commit-config.yaml`

## Key Files

- `./main.go` - CLI entry point with flag parsing
- `internal/api/server.go` - HTTP server and API routes
- `internal/config/config.go` - Configuration loading (JSON file + env vars)
- `internal/parser/dmarc.go` - DMARC XML parsing logic
- `internal/imap/client.go` - IMAP email fetching
- `internal/storage/common.go` - Shared storage logic and SQL queries
- `internal/storage/sqlite_no_cgo.go` - Pure Go SQLite implementation (modernc.org/sqlite)
- `internal/storage/sqlite_cgo.go` - CGO SQLite implementation (mattn/go-sqlite3)
- `src/App.vue` - Main Vue.js dashboard component
- `src/components/DnsGenerator.vue` - DMARC DNS record generator component

## API Endpoints

- `GET /api/statistics` - Dashboard statistics
- `GET /api/reports` - List reports (paginated)
- `GET /api/reports/:id` - Single report details
- `GET /api/top-sources` - Top sending source IPs

## Frontend Features

The Vue.js dashboard includes:

- **Dashboard Statistics** - Overview of DMARC report metrics
- **Reports List** - Paginated view of all DMARC reports
- **Report Details** - Detailed view of individual reports
- **Top Sources** - Visualization of top sending source IPs
- **DMARC DNS Generator** - Interactive tool to generate DMARC DNS TXT records with configurable policies (none/quarantine/reject), subdomain policies, reporting options (rua/ruf), and alignment modes

## Configuration

Config via JSON file or environment variables (using caarlos0/env):

- IMAP settings (host, port, username, password, mailbox, TLS)
- Database path
- Server host/port

Example: `config.example.json`

## CI/CD

- GitHub Actions workflow in `.github/workflows/ci.yml`
- golangci-lint for Go linting
- Docker build with cosign signing
- Release automation via release-please and goreleaser

## Roadmap

See `ROADMAP.md` for the comprehensive product roadmap including:

- Phase 1: Delightful Defaults (dark mode, DNS generator, health score, exports)
- Phase 2: Proactive Intelligence (alerting, trends, GeoIP maps, DNS validator)
- Phase 3: Enterprise Ready (auth, multi-org, RBAC, Prometheus metrics)
- Phase 4: AI-Powered Security (AI assistant, anomaly detection, forensic reports)

## Contributing

See `CONTRIBUTING.md` for development setup and contribution guidelines. Key areas for contribution:

- Forensic Reports (RUF support)
- OAuth2 for IMAP authentication
- CSV/JSON export functionality
- Email alerts for compliance issues
- Historical trend analysis
- Test coverage improvements
