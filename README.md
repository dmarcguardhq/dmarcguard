# Parse DMARC

[![CI](https://img.shields.io/github/actions/workflow/status/meysam81/parse-dmarc/ci.yml?branch=main&label=CI&logo=githubactions&logoColor=white&style=flat-square)](https://github.com/meysam81/parse-dmarc/actions/workflows/ci.yml)
[![License](https://img.shields.io/github/license/meysam81/parse-dmarc?style=flat-square)](https://github.com/meysam81/parse-dmarc/blob/main/LICENSE)
[![GitHub release](https://img.shields.io/github/v/release/meysam81/parse-dmarc?style=flat-square&logo=github)](https://github.com/meysam81/parse-dmarc/releases)
[![GitHub Stars](https://img.shields.io/github/stars/meysam81/parse-dmarc?style=flat-square&logo=github)](https://github.com/meysam81/parse-dmarc/stargazers)
[![GitHub Issues](https://img.shields.io/github/issues/meysam81/parse-dmarc?style=flat-square&logo=github)](https://github.com/meysam81/parse-dmarc/issues)
[![Go Report Card](https://goreportcard.com/badge/github.com/meysam81/parse-dmarc?style=flat-square)](https://goreportcard.com/report/github.com/meysam81/parse-dmarc)

[![Made with Go](https://img.shields.io/badge/Made%20with-Go-1f425f?style=flat-square&logo=go&logoColor=white)](https://go.dev)
[![Made with Vue.js](https://img.shields.io/badge/Made%20with-Vue.js-4FC08D?style=flat-square&logo=vue.js&logoColor=white)](https://vuejs.org)
[![Docker Hub](https://img.shields.io/badge/Docker%20Hub-meysam81%2Fparse--dmarc-2496ED?style=flat-square&logo=docker&logoColor=white)](https://hub.docker.com/r/meysam81/parse-dmarc)
[![Docker Pulls](https://img.shields.io/docker/pulls/meysam81/parse-dmarc?style=flat-square&logo=docker&logoColor=white)](https://hub.docker.com/r/meysam81/parse-dmarc)
[![Docker Image Size (tag)](https://img.shields.io/docker/image-size/meysam81/parse-dmarc/v1?style=flat-square&logo=docker&logoColor=white)](https://hub.docker.com/r/meysam81/parse-dmarc)

[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-FE5196?logo=conventionalcommits&logoColor=white&style=flat-square)](https://www.conventionalcommits.org)
[![Renovate](https://img.shields.io/badge/renovate-enabled-1f8b4c?logo=renovatebot&logoColor=white&style=flat-square)](https://developer.mend.io/github/meysam81/parse-dmarc)

**Monitor who's sending email on behalf of your domain. Catch spoofing. Stop phishing.**

[![Parse DMARC](./assets/social-preview.png)](https://github.com/meysam81/parse-dmarc)

## Deploy Your Own Instance

Deploy Parse DMARC to your favorite cloud provider with one click:

### Platform as a Service (PaaS)

| Provider       | Deploy                                                                                                                                                                             | Notes                                                     |
| -------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------- |
| **Railway**    | [![Deploy on Railway](https://railway.com/button.svg)](https://railway.com/deploy/4kqQ_I?referralCode=meysam)                                                                      | Recommended for beginners                                 |
| **Render**     | [![Deploy to Render](https://render.com/images/deploy-to-render-button.svg)](https://render.com/deploy?repo=https://github.com/meysam81/parse-dmarc)                               | Free tier available                                       |
| **Koyeb**      | [![Deploy to Koyeb](https://www.koyeb.com/static/images/deploy/button.svg)][koyeb-1click]                                                                                          | Global edge deployment. Manually mount `/data` as volume. |
| **Zeabur**     | [![Deploy on Zeabur](https://zeabur.com/button.svg)](https://zeabur.com/templates/YB3TS7?referralCode=meysam)                                                                      | Asia-Pacific optimized                                    |
| **Northflank** | [![Deploy to Northflank](https://assets.northflank.com/deploy_to_northflank_smm_36700fb050.svg)](https://app.northflank.com/s/account/templates/new?data=693e394eb41e1e64db65187e) | Developer-focused                                         |

### Self-Hosted

| Provider     | Deploy                                                                                                                                              | Notes                           |
| ------------ | --------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------- |
| **CapRover** | [![Deploy to CapRover](https://img.shields.io/badge/Deploy-CapRover-0072CE?style=for-the-badge&logo=docker)](./deploy/captain-definition)           | Self-hosted PaaS                |
| **Coolify**  | [![Deploy to Coolify](https://img.shields.io/badge/Deploy-Coolify-6B46C1?style=for-the-badge&logo=docker)](./deploy/coolify.yaml)                   | Open-source Heroku alternative  |
| **Dokploy**  | [![Deploy to Dokploy](https://img.shields.io/badge/Deploy-Dokploy-00B4D8?style=for-the-badge&logo=docker)](./deploy/dokploy/)                       | Self-hosted deployment platform |
| **Docker**   | [![Docker](https://img.shields.io/badge/Docker-Pull%20Image-2496ED?style=for-the-badge&logo=docker)](https://hub.docker.com/r/meysam81/parse-dmarc) | Run anywhere                    |

### Infrastructure

| Provider                 | Deploy                                                                                                                                             | Notes                |
| ------------------------ | -------------------------------------------------------------------------------------------------------------------------------------------------- | -------------------- |
| **DigitalOcean Droplet** | [![Deploy to DigitalOcean](https://img.shields.io/badge/Deploy-DigitalOcean-0080FF?style=for-the-badge&logo=digitalocean)](./deploy/digitalocean/) | VM with Packer image |

> **Note**: All deployments require IMAP credentials. See [Configuration](#configuration-options) for details on setting up Gmail, Outlook, or other email providers.

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
- 🔐 Optional GitHub OAuth login for the dashboard ([details](#dashboard-authentication))
- 🌙 Dark mode support

## Installation

### Homebrew (macOS/Linux)

```bash
brew tap meysam81/tap
brew install parse-dmarc
```

### Docker

```bash
docker pull meysam81/parse-dmarc
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

**Run the container:**

```bash
docker run -d \
  --name parse-dmarc \
  -p 8080:8080 \
  -e IMAP_HOST=imap.gmail.com \
  -e IMAP_PORT=993 \
  -e IMAP_USERNAME=your-email@gmail.com \
  -e IMAP_PASSWORD=your-app-password \
  -v parse-dmarc:/data \
  meysam81/parse-dmarc
```

**For Gmail users:** You'll need an [App Password](https://support.google.com/accounts/answer/185833), not your regular Gmail password.

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

## Dashboard Authentication

The dashboard runs unauthenticated by default — fine for `localhost`, but if
you expose it to the internet you'll want login. Parse DMARC ships with
**GitHub OAuth** out of the box.

### Step 1 — Create a GitHub OAuth App

Go to <https://github.com/settings/developers> → **OAuth Apps** → **New OAuth App**:

| Field | Value |
| --- | --- |
| **Application name** | `dmarcguard` (anything — only you see it) |
| **Homepage URL** | Your dashboard URL, e.g. `https://dmarc.example.com` (or `http://localhost:8080` for local testing) |
| **Authorization callback URL** | Same as above with `/auth/callback` appended, e.g. `https://dmarc.example.com/auth/callback` |
| **Enable Device Flow** | ❌ Leave unchecked |

Click **Register application**. On the next screen:

1. Copy the **Client ID** (visible immediately).
2. Click **Generate a new client secret** → copy the secret (only shown once).

> The callback URL above must match `auth.redirect_url` in your config **exactly** — including scheme (`http`/`https`), host, port, and path. A mismatch produces GitHub's `redirect_uri_mismatch` error.

### Step 2 — Generate a session secret

```bash
./parse-dmarc --gen-session-secret
# nX2Fc0PbLjYOvGPrO3MsZ6o+To2MHi3hJp48czUnH+U=
```

Or in Docker:

```bash
docker compose run --rm parse-dmarc --gen-session-secret
```

This is a 32-byte random key used to sign session cookies (HMAC-SHA256). Treat it like any other secret — anyone who has it can mint valid sessions for any allowlisted user.

### Step 3 — Add the `auth` block to `config.json`

```json
"auth": {
  "enabled": true,
  "client_id": "Iv1.xxxxxxxxxxxxxxxx",
  "client_secret": "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
  "redirect_url": "https://dmarc.example.com/auth/callback",
  "session_secret": "nX2Fc0PbLjYOvGPrO3MsZ6o+To2MHi3hJp48czUnH+U=",
  "allowed_users": ["sebykrueger"],
  "allowed_emails": ["seb@example.com"],
  "session_ttl_days": 7
}
```

**Allowlist semantics:**

- `allowed_users` matches the GitHub username (e.g. `sebykrueger`).
- `allowed_emails` matches the verified primary email on the GitHub account.
- A user is granted access if **either** list contains their identity (OR, not AND). Case-insensitive on both sides.
- At least one entry across the two lists is required — the daemon refuses to start with `auth.enabled=true` and an empty allowlist (otherwise nobody could log in).

**Optional fields:**

- `session_ttl_days` — how long a login lasts. Default: 7 days.
- The `Secure` cookie flag is set automatically when `redirect_url` starts with `https://`.

### Step 4 — Restart the daemon

```bash
docker compose restart parse-dmarc
```

You should see this log line confirming auth is on:

```
INF dashboard authentication enabled (github) allowed_emails=1 allowed_users=1
```

### What changes after auth is enabled

| Path | Behavior |
| --- | --- |
| `/` and the dashboard SPA | 303 redirect to `/auth/login` if no valid session cookie |
| `/api/*` | `401 application/json` with `{"error":"unauthorized","login_url":"/auth/login"}` for programmatic callers |
| `/auth/login` | Redirects to GitHub OAuth |
| `/auth/callback` | Receives the code, validates state (CSRF), checks the allowlist, sets the session cookie |
| `/auth/logout` | Clears the session cookie and redirects to login |
| `/metrics` | **Still open** — Prometheus scrapers don't need to authenticate |

### Environment variable equivalents

Every field has an env var, useful for Docker secrets / Vault / k8s secrets:

```bash
AUTH_ENABLED=true
AUTH_CLIENT_ID=Iv1.xxxxxxxxxxxxxxxx
AUTH_CLIENT_SECRET=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
AUTH_REDIRECT_URL=https://dmarc.example.com/auth/callback
AUTH_SESSION_SECRET=nX2Fc0PbLjYOvGPrO3MsZ6o+To2MHi3hJp48czUnH+U=
AUTH_ALLOWED_USERS=sebykrueger,alice
AUTH_ALLOWED_EMAILS=seb@example.com
AUTH_SESSION_TTL_DAYS=7
```

### Common gotchas

- **`redirect_uri_mismatch` from GitHub** — the OAuth App's callback URL doesn't match `auth.redirect_url` to the byte. Trailing slash counts.
- **403 "not authorized"** after a successful GitHub login — your username/email isn't in the allowlist. Either add it, or switch to a different GitHub account.
- **Removing someone from the allowlist doesn't kick them out immediately** — sessions are cookie-based and last `session_ttl_days`. To force-evict, rotate `auth.session_secret` (invalidates all sessions).
- **Cookies aren't `Secure` on `http://`** — fine for local, dangerous in production. Always use HTTPS for any internet-facing deployment (a reverse proxy with Caddy / Traefik / nginx-acme handles this in one line).

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
docker run -e PARSE_DMARC_METRICS=false meysam81/parse-dmarc
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
    image: meysam81/parse-dmarc
    ports:
      - "8080:8080"
    volumes:
      - ./config.json:/app/config.json
      - ./data:/data

  prometheus:
    image: prom/prometheus
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
    command:
      - "--config.file=/etc/prometheus/prometheus.yml"

  grafana:
    image: grafana/grafana
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

[koyeb-1click]: https://app.koyeb.com/deploy?name=parse-dmarc&type=docker&image=docker.io%2Fmeysam81%2Fparse-dmarc%3Alatest&regions=fra&env%5BDATABASE_PATH%5D=%2Fdata%2Fdb.sqlite&env%5BIMAP_HOST%5D=&env%5BIMAP_MAILBOX%5D=INBOX&env%5BIMAP_PASSWORD%5D=&env%5BIMAP_PORT%5D=993&env%5BIMAP_USERNAME%5D=&env%5BIMAP_USE_TLS%5D=true&env%5BSERVER_PORT%5D=8080&ports=8080%3Bhttp%3B%2F&hc_protocol%5B8080%5D=http&hc_grace_period%5B8080%5D=5&hc_interval%5B8080%5D=30&hc_restart_limit%5B8080%5D=3&hc_timeout%5B8080%5D=5&hc_path%5B8080%5D=%2Fapi%2Fstatistics&hc_method%5B8080%5D=get
