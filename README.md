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
  -v $(pwd)/data:/data \
  ghcr.io/meysam81/parse-dmarc:latest
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
