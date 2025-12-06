# Parse DMARC

[![License](https://img.shields.io/github/license/meysam81/parse-dmarc?style=flat-square)](https://github.com/meysam81/parse-dmarc/blob/main/LICENSE)
[![GitHub release](https://img.shields.io/github/v/release/meysam81/parse-dmarc?style=flat-square)](https://github.com/meysam81/parse-dmarc/releases)
[![GitHub Stars](https://img.shields.io/github/stars/meysam81/parse-dmarc?style=flat-square&logo=github)](https://github.com/meysam81/parse-dmarc/stargazers)
[![GitHub Issues](https://img.shields.io/github/issues/meysam81/parse-dmarc?style=flat-square&logo=github)](https://github.com/meysam81/parse-dmarc/issues)
[![Go Report Card](https://goreportcard.com/badge/github.com/meysam81/parse-dmarc?style=flat-square)](https://goreportcard.com/report/github.com/meysam81/parse-dmarc)
[![Made with Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg?style=flat-square&logo=go)](https://go.dev)
[![Made with Vue.js](https://img.shields.io/badge/Made%20with-Vue.js-4FC08D.svg?style=flat-square&logo=vue.js)](https://vuejs.org)
[![GHCR](https://img.shields.io/badge/GHCR-ghcr.io%2Fmeysam81%2Fparse--dmarc-blue?style=flat-square&logo=github)](https://github.com/meysam81/parse-dmarc/pkgs/container/parse-dmarc)

**Monitor who's sending email on behalf of your domain. Catch spoofing. Stop phishing.**

[![Dashboard Screenshot](./assets/demo.png)](https://github.com/meysam81/parse-dmarc)

## Why Do I Need This?

**DMARC** (Domain-based Message Authentication, Reporting & Conformance) helps protect your domain from email spoofing and phishing. When you enable DMARC on your domain, email providers like Gmail, Outlook, and Yahoo send you **aggregate reports** showing:

- Who's sending email claiming to be from your domain
- Which emails passed or failed authentication (SPF/DKIM)
- How many emails were sent, and from which IP addresses
- Whether malicious actors are trying to impersonate your domain

**The Problem:** These reports arrive as compressed XML attachments in your inbox - nearly impossible to read or analyze manually.

**The Solution:** Parse DMARC automatically fetches these reports from your inbox, parses them, and displays everything in a beautiful dashboard. All in a single 14MB Docker image.

## Features

- 📧 Auto-fetches reports from any IMAP inbox (Gmail, Outlook, etc.)
- 📊 Beautiful dashboard with real-time statistics
- 🔍 See exactly who's sending email as your domain
- 🔧 Built-in DNS record generator for easy DMARC setup
- 📦 Single binary - no databases to install, no complex setup
- 🚀 Tiny 14MB Docker image
- 🔒 Secure TLS support
- 🌙 Dark mode support

## Installation

### Homebrew (macOS/Linux)

```bash
brew install meysam81/tap/parse-dmarc
```

Or tap first, then install:

```bash
brew tap meysam81/tap
brew install parse-dmarc
```

### Docker

```bash
docker pull ghcr.io/meysam81/parse-dmarc:latest
```

### Binary Downloads

Download pre-built binaries from the [Releases page](https://github.com/meysam81/parse-dmarc/releases).

## Quick Start

### Step 1: Set Up DNS to Receive DMARC Reports

**This is the most important step!** Without this, you won't receive any reports to analyze.

Add a DMARC TXT record to your domain's DNS:

```
Name: _dmarc.yourdomain.com
Type: TXT
Value: v=DMARC1; p=none; rua=mailto:dmarc@yourdomain.com
```

**What this means:**

- `p=none` - Monitor only (don't block emails yet)
- `rua=mailto:dmarc@yourdomain.com` - Send aggregate reports to this email address

**Important:** Replace `dmarc@yourdomain.com` with an actual email inbox you control. This is where Gmail, Outlook, Yahoo, etc. will send your DMARC reports.

**DNS Examples:**

- **Cloudflare:** DNS > Add record > Type: TXT, Name: `_dmarc`, Content: `v=DMARC1; p=none; rua=mailto:dmarc@yourdomain.com`
- **Google Domains:** DNS > Custom records > TXT, Name: `_dmarc`, Data: `v=DMARC1; p=none; rua=mailto:dmarc@yourdomain.com`
- **AWS Route53:** Create record > Type: TXT, Name: `_dmarc.yourdomain.com`, Value: `"v=DMARC1; p=none; rua=mailto:dmarc@yourdomain.com"`

Reports typically start arriving within 24-48 hours.

### Step 2: Run Parse DMARC with Docker

**Create a configuration file:**

```bash
mkdir -p data
cat > config.json <<EOF
{
  "imap": {
    "host": "imap.gmail.com",
    "port": 993,
    "username": "dmarc@yourdomain.com",
    "password": "your-app-password",
    "mailbox": "INBOX",
    "use_tls": true
  },
  "database": {
    "path": "/data/db.sqlite"
  },
  "server": {
    "port": 8080,
    "host": "0.0.0.0"
  }
}
EOF
```

**For Gmail users:** You'll need an [App Password](https://support.google.com/accounts/answer/185833), not your regular Gmail password.

**Run the container:**

```bash
docker run -d \
  --name parse-dmarc \
  -p 8080:8080 \
  -v $(pwd)/config.json:/app/config.json \
  -v parse-dmarc-data:/data \
  ghcr.io/meysam81/parse-dmarc:v1.2.0
```

**Access the dashboard:** Open `http://localhost:8080` in your browser.

## What You'll See

Once DMARC reports start arriving and Parse DMARC processes them, your dashboard will show:

- **Total messages** analyzed across all reports
- **DMARC compliance rate** (SPF/DKIM pass rates)
- **Top sending sources** (IP addresses and organizations sending as your domain)
- **Authentication results** (which emails passed/failed SPF and DKIM)
- **Policy actions** (how receiving servers handled your email)

This helps you:

- Verify your legitimate email services are properly configured
- Detect unauthorized use of your domain
- Gradually move from monitoring (`p=none`) to enforcement (`p=quarantine` or `p=reject`)

## Configuration Options

### IMAP Settings for Common Providers

**Gmail:**

```json
{
  "host": "imap.gmail.com",
  "port": 993,
  "username": "your-email@gmail.com",
  "password": "your-app-password",
  "use_tls": true
}
```

Requires [App Password](https://support.google.com/accounts/answer/185833)

**Outlook/Office 365:**

```json
{
  "host": "outlook.office365.com",
  "port": 993,
  "username": "your-email@outlook.com",
  "password": "your-password",
  "use_tls": true
}
```

**Generic IMAP:**
Most providers use port `993` with TLS. Check your provider's documentation.

### Command Line Options

```bash
# Fetch once and exit (useful for cron jobs)
docker exec parse-dmarc ./parse-dmarc -fetch-once

# Serve dashboard only (no fetching)
docker exec parse-dmarc ./parse-dmarc -serve-only

# Custom fetch interval (in seconds, default 300)
docker exec parse-dmarc ./parse-dmarc -fetch-interval=600
```

## Frequently Asked Questions

**Q: I'm not receiving any reports. What's wrong?**

A: Check these things in order:

1. Did you add the `_dmarc` TXT record to your DNS? (Use a DNS checker like `dig _dmarc.yourdomain.com TXT`)
2. Wait 24-48 hours - reports aren't instant
3. Is your domain sending/receiving email? No email = no reports
4. Check your IMAP credentials are correct in `config.json`

**Q: Do I need SPF and DKIM set up first?**

A: No! DMARC reports will show you whether SPF and DKIM are passing or failing, which helps you configure them correctly.

**Q: What should my DMARC policy be?**

A: Start with `p=none` (monitoring only). After reviewing reports and fixing any issues, gradually move to `p=quarantine` and then `p=reject`.

**Q: How much email traffic do I need?**

A: Any amount works. Even small domains with a few emails per day will receive useful reports.

**Q: Can I use a Gmail account to receive reports?**

A: Yes! Create a dedicated Gmail like `dmarc@yourdomain.com`, forward it to your personal Gmail if needed, and use Gmail's IMAP settings.

## Advanced

### Building from Source

```bash
git clone https://github.com/meysam81/parse-dmarc.git
cd parse-dmarc
just install-deps
just build
./bin/parse-dmarc -config=config.json
```

### Docker Compose

See [`compose.yml`](./compose.yml) for Docker Compose configuration.

### API Endpoints

- `GET /api/statistics` - Dashboard statistics
- `GET /api/reports` - List of reports (paginated)
- `GET /api/reports/:id` - Detailed report view
- `GET /api/top-sources` - Top sending source IPs
- `GET /metrics` - Prometheus metrics endpoint

## Prometheus Metrics & Grafana Integration

Parse DMARC includes production-ready Prometheus metrics for monitoring and alerting. Metrics are enabled by default and exposed at `/metrics`.

### Available Metrics

#### Build Information

| Metric                   | Type  | Description                                     |
| ------------------------ | ----- | ----------------------------------------------- |
| `parse_dmarc_build_info` | Gauge | Build information (version, commit, build_date) |

#### Report Processing

| Metric                                             | Type      | Description                                 |
| -------------------------------------------------- | --------- | ------------------------------------------- |
| `parse_dmarc_reports_fetched_total`                | Counter   | Total DMARC report emails fetched from IMAP |
| `parse_dmarc_reports_parsed_total`                 | Counter   | Total DMARC reports successfully parsed     |
| `parse_dmarc_reports_stored_total`                 | Counter   | Total DMARC reports stored in database      |
| `parse_dmarc_reports_parse_errors_total`           | Counter   | Total parse errors                          |
| `parse_dmarc_reports_store_errors_total`           | Counter   | Total storage errors                        |
| `parse_dmarc_reports_attachments_total`            | Counter   | Total attachments processed                 |
| `parse_dmarc_reports_fetch_duration_seconds`       | Histogram | Duration of fetch operations                |
| `parse_dmarc_reports_last_fetch_timestamp_seconds` | Gauge     | Unix timestamp of last successful fetch     |
| `parse_dmarc_reports_fetch_cycles_total`           | Counter   | Total fetch cycles executed                 |
| `parse_dmarc_reports_fetch_errors_total`           | Counter   | Total fetch cycle errors                    |

#### IMAP Connection

| Metric                                         | Type      | Labels | Description                              |
| ---------------------------------------------- | --------- | ------ | ---------------------------------------- |
| `parse_dmarc_imap_connections_total`           | Counter   | status | IMAP connection attempts (success/error) |
| `parse_dmarc_imap_connection_duration_seconds` | Histogram |        | IMAP connection establishment duration   |

#### DMARC Statistics

| Metric                                       | Type  | Description                       |
| -------------------------------------------- | ----- | --------------------------------- |
| `parse_dmarc_dmarc_reports_total`            | Gauge | Total reports in database         |
| `parse_dmarc_dmarc_messages_total`           | Gauge | Total messages across all reports |
| `parse_dmarc_dmarc_compliant_messages_total` | Gauge | Total DMARC-compliant messages    |
| `parse_dmarc_dmarc_compliance_rate`          | Gauge | Overall compliance rate (0-100)   |
| `parse_dmarc_dmarc_unique_source_ips`        | Gauge | Number of unique source IPs       |
| `parse_dmarc_dmarc_unique_domains`           | Gauge | Number of unique domains          |

#### Per-Domain/Org Metrics

| Metric                                        | Type  | Labels      | Description                  |
| --------------------------------------------- | ----- | ----------- | ---------------------------- |
| `parse_dmarc_dmarc_messages_by_domain`        | Gauge | domain      | Messages per domain          |
| `parse_dmarc_dmarc_compliance_rate_by_domain` | Gauge | domain      | Compliance rate per domain   |
| `parse_dmarc_dmarc_reports_by_org`            | Gauge | org_name    | Reports per organization     |
| `parse_dmarc_dmarc_messages_by_disposition`   | Gauge | disposition | Messages by disposition type |

#### Authentication Results

| Metric                           | Type  | Labels | Description                       |
| -------------------------------- | ----- | ------ | --------------------------------- |
| `parse_dmarc_dmarc_spf_results`  | Gauge | result | SPF authentication result counts  |
| `parse_dmarc_dmarc_dkim_results` | Gauge | result | DKIM authentication result counts |

#### HTTP Server

| Metric                                      | Type      | Labels               | Description                |
| ------------------------------------------- | --------- | -------------------- | -------------------------- |
| `parse_dmarc_http_requests_total`           | Counter   | method, path, status | Total HTTP requests        |
| `parse_dmarc_http_request_duration_seconds` | Histogram | method, path         | HTTP request duration      |
| `parse_dmarc_http_requests_in_flight`       | Gauge     |                      | Current in-flight requests |

#### Go Runtime (Built-in)

Standard Go runtime metrics are also exposed:

- `go_goroutines` - Number of goroutines
- `go_memstats_*` - Memory statistics
- `go_gc_*` - Garbage collection metrics
- `process_*` - Process metrics (CPU, memory, file descriptors)

### Disabling Metrics

To disable the metrics endpoint:

```bash
# Command line
./parse-dmarc --metrics=false

# Environment variable
export PARSE_DMARC_METRICS=false

# Docker
docker run -e PARSE_DMARC_METRICS=false ghcr.io/meysam81/parse-dmarc:latest
```

### Prometheus Configuration

Add Parse DMARC to your `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: "parse-dmarc"
    static_configs:
      - targets: ["parse-dmarc:8080"]
    scrape_interval: 30s
    metrics_path: /metrics
```

For Kubernetes with ServiceMonitor (Prometheus Operator):

```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: parse-dmarc
  labels:
    app: parse-dmarc
spec:
  selector:
    matchLabels:
      app: parse-dmarc
  endpoints:
    - port: http
      path: /metrics
      interval: 30s
```

### Grafana Dashboard

A production-ready Grafana dashboard is included in `grafana/dashboard.json`.

#### Import Manually

1. In Grafana, go to **Dashboards** > **Import**
2. Upload `grafana/dashboard.json` or paste its contents
3. Select your Prometheus datasource
4. Click **Import**

#### Provision Automatically (Recommended for Production)

```bash
# Copy dashboard to Grafana dashboards directory
cp grafana/dashboard.json /var/lib/grafana/dashboards/parse-dmarc/

# Copy provisioning config
cp grafana/provisioning.yaml /etc/grafana/provisioning/dashboards/parse-dmarc.yaml

# Restart Grafana or wait for it to pick up changes
systemctl restart grafana-server
```

#### Dashboard Variables

| Variable     | Purpose                        |
| ------------ | ------------------------------ |
| `datasource` | Prometheus datasource to query |
| `job`        | Filter by Prometheus job label |
| `instance`   | Filter by instance(s)          |
| `domain`     | Filter by monitored domain(s)  |

#### Dashboard Sections

| Section                            | What It Shows                                                             |
| ---------------------------------- | ------------------------------------------------------------------------- |
| **Overview - Golden Signals**      | Compliance rate, total messages, reports count, time since last fetch     |
| **DMARC Authentication Results**   | SPF/DKIM pass rates, disposition breakdown, per-domain compliance         |
| **Report Sources & Organizations** | Top reporting organizations (Google, Microsoft, etc.), messages by domain |
| **IMAP & Fetch Operations**        | Connection health, fetch cycle monitoring, latency heatmaps               |
| **Error Tracking**                 | Parse errors, storage errors, fetch failures                              |
| **HTTP Server**                    | Request rates, latency percentiles, error rates                           |
| **Go Runtime**                     | Goroutines, memory usage, GC stats, CPU usage                             |

#### Example Grafana Panels

**Compliance Rate Gauge:**

```promql
parse_dmarc_dmarc_compliance_rate
```

**Messages Over Time:**

```promql
rate(parse_dmarc_dmarc_messages_total[5m])
```

**Compliance Rate by Domain:**

```promql
parse_dmarc_dmarc_compliance_rate_by_domain
```

**SPF/DKIM Pass Rate:**

```promql
# SPF Pass Rate
parse_dmarc_dmarc_spf_results{result="pass"} / ignoring(result) sum(parse_dmarc_dmarc_spf_results) * 100

# DKIM Pass Rate
parse_dmarc_dmarc_dkim_results{result="pass"} / ignoring(result) sum(parse_dmarc_dmarc_dkim_results) * 100
```

**Fetch Success Rate:**

```promql
1 - (rate(parse_dmarc_reports_fetch_errors_total[1h]) / rate(parse_dmarc_reports_fetch_cycles_total[1h]))
```

**IMAP Connection Health:**

```promql
rate(parse_dmarc_imap_connections_total{status="success"}[5m]) /
(rate(parse_dmarc_imap_connections_total{status="success"}[5m]) + rate(parse_dmarc_imap_connections_total{status="error"}[5m]))
```

**HTTP Request Latency (p95):**

```promql
histogram_quantile(0.95, rate(parse_dmarc_http_request_duration_seconds_bucket[5m]))
```

**Reports by Organization:**

```promql
topk(10, parse_dmarc_dmarc_reports_by_org)
```

#### Alerting Rules

Example Prometheus alerting rules:

```yaml
groups:
  - name: parse-dmarc
    rules:
      - alert: DMARCComplianceLow
        expr: parse_dmarc_dmarc_compliance_rate < 90
        for: 1h
        labels:
          severity: warning
        annotations:
          summary: "DMARC compliance rate is below 90%"
          description: "Current compliance rate: {{ $value }}%"

      - alert: DMARCFetchFailures
        expr: rate(parse_dmarc_reports_fetch_errors_total[15m]) > 0
        for: 30m
        labels:
          severity: critical
        annotations:
          summary: "Parse DMARC fetch failures detected"
          description: "IMAP fetch operations are failing"

      - alert: IMAPConnectionErrors
        expr: rate(parse_dmarc_imap_connections_total{status="error"}[5m]) > 0
        for: 10m
        labels:
          severity: warning
        annotations:
          summary: "IMAP connection errors detected"
          description: "Check IMAP credentials and server connectivity"

      - alert: NoRecentFetch
        expr: time() - parse_dmarc_reports_last_fetch_timestamp_seconds > 600
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "No recent DMARC report fetch"
          description: "Last fetch was {{ humanizeDuration $value }} ago"
```

### Docker Compose with Prometheus & Grafana

Complete monitoring stack:

```yaml
version: "3.8"

services:
  parse-dmarc:
    image: ghcr.io/meysam81/parse-dmarc:latest
    ports:
      - "8080:8080"
    volumes:
      - ./config.json:/app/config.json
      - ./data:/data

  prometheus:
    image: prom/prometheus:latest
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
    command:
      - "--config.file=/etc/prometheus/prometheus.yml"

  grafana:
    image: grafana/grafana:latest
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
    volumes:
      - grafana-data:/var/lib/grafana

volumes:
  grafana-data:
```

With `prometheus.yml`:

```yaml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: "parse-dmarc"
    static_configs:
      - targets: ["parse-dmarc:8080"]
```

Access:

- Parse DMARC Dashboard: http://localhost:8080
- Prometheus: http://localhost:9090
- Grafana: http://localhost:3000 (admin/admin)

### Why Parse DMARC vs ParseDMARC?

This project is inspired by [ParseDMARC](https://github.com/domainaware/parsedmarc) but built for simplicity:

- **Single 14MB binary** vs Python + Elasticsearch + Kibana stack
- **Built-in dashboard** vs external visualization tools
- **SQLite** vs Elasticsearch (no JVM required)
- **Zero dependencies** vs complex setup

## Contributing

Issues and pull requests are welcome! Please check the [issues page](https://github.com/meysam81/parse-dmarc/issues).

## License

Apache-2.0 - see [LICENSE](LICENSE) for details.

---

**Found this useful? Star the repo!** ⭐
