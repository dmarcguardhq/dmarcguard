# Parse DMARC Product Roadmap

> **Vision**: Make email authentication accessible to everyone—from solo developers to enterprises—with a tool so intuitive and powerful that protecting your domain becomes effortless.

This roadmap outlines our strategic direction for Parse DMARC. It's designed to transform the project from a solid DMARC parser into the definitive open-source email authentication platform.

---

## Table of Contents

- [Our Users](#our-users)
- [Design Principles](#design-principles)
- [Roadmap Phases](#roadmap-phases)
  - [Phase 1: Delightful Defaults](#phase-1-delightful-defaults)
  - [Phase 2: Proactive Intelligence](#phase-2-proactive-intelligence)
  - [Phase 3: Enterprise Ready](#phase-3-enterprise-ready)
  - [Phase 4: AI-Powered Security](#phase-4-ai-powered-security)
- [Feature Details](#feature-details)
- [Future Exploration](#future-exploration)
- [Contributing](#contributing)

---

## Our Users

Understanding who we're building for shapes every decision. Parse DMARC serves several distinct personas:

### The Solo Developer

- **Context**: Runs multiple side projects, each with its own domain
- **Pain Point**: No time to learn email authentication intricacies; just wants things to work
- **Need**: Set-and-forget solution with instant visibility when something goes wrong
- **Success Metric**: Protected in under 5 minutes, never thinks about DMARC again until needed

### The Startup Engineer

- **Context**: Building a SaaS product, sending transactional and marketing emails
- **Pain Point**: Email deliverability issues affecting customer experience; no budget for enterprise tools
- **Need**: Comprehensive monitoring without the Elasticsearch/Kibana complexity
- **Success Metric**: Full visibility into email authentication, catches issues before customers complain

### The Security Engineer

- **Context**: Protecting organization's domain from spoofing and phishing attacks
- **Pain Point**: Aggregate reports are unreadable; can't detect attacks in real-time
- **Need**: Threat intelligence, alerting, forensic analysis capabilities
- **Success Metric**: Detects and responds to spoofing attempts within minutes

### The MSP/Agency

- **Context**: Managing email security for dozens or hundreds of client domains
- **Pain Point**: No unified view across clients; manual report checking doesn't scale
- **Need**: Multi-tenant dashboard, white-labeling, automated reporting
- **Success Metric**: Manages 100+ domains from a single pane of glass

### The Enterprise IT Team

- **Context**: Large organization with complex email infrastructure
- **Pain Point**: Compliance requirements, audit trails, integration with existing security stack
- **Need**: SSO, RBAC, SIEM integration, comprehensive audit logging
- **Success Metric**: Passes security audits, integrates with existing tooling

---

## Design Principles

Every feature we build adheres to these principles:

### 1. Zero-to-Protected in Minutes

The time from "I heard about DMARC" to "my domain is monitored" should be measured in minutes, not hours. Every friction point is a failure.

### 2. Actionable, Not Just Informational

Data without context is noise. Every metric, alert, and visualization should answer "so what?" and ideally "what should I do about it?"

### 3. Complexity is Optional

Simple use cases should feel simple. Advanced features should be discoverable but never in the way. A solo developer and an enterprise security team should both feel at home.

### 4. Works Where You Work

Parse DMARC should integrate into existing workflows—Slack, Prometheus, CI/CD pipelines—rather than demanding users change how they work.

### 5. Trust Through Transparency

As a security tool, we must be transparent about what we do, how we handle data, and what our limitations are. Open source isn't just our license; it's our philosophy.

---

## Roadmap Phases

### Phase 1: Delightful Defaults

**Theme**: Polish the foundation, make first impressions unforgettable

This phase focuses on quick wins that dramatically improve user experience without major architectural changes. Every feature here should make users think "why doesn't every tool do this?"

| Feature                                                       | Priority | Status      |
| ------------------------------------------------------------- | -------- | ----------- |
| [Dark Mode](#11-dark-mode)                                    | High     | Planned     |
| [DNS Record Generator](#12-dns-record-generator)              | High     | ✅ Complete |
| [Compliance Health Score](#13-compliance-health-score)        | High     | Planned     |
| [Data Export (CSV/JSON)](#14-data-export)                     | High     | Planned     |
| [Keyboard Shortcuts](#15-keyboard-shortcuts)                  | Medium   | Planned     |
| [ASN/Organization Enrichment](#16-asnorganization-enrichment) | Medium   | Planned     |
| [Compliance Badge](#17-compliance-badge)                      | Medium   | Planned     |
| [Manual Report Upload](#18-manual-report-upload)              | Medium   | Planned     |

---

### Phase 2: Proactive Intelligence

**Theme**: Transform from passive monitoring to active protection

This phase adds the intelligence layer that turns Parse DMARC from a report viewer into a security tool. Users should feel protected, not just informed.

| Feature                                                    | Priority | Status  |
| ---------------------------------------------------------- | -------- | ------- |
| [Smart Alerting System](#21-smart-alerting-system)         | High     | Planned |
| [Historical Trend Charts](#22-historical-trend-charts)     | High     | Planned |
| [Interactive GeoIP Map](#23-interactive-geoip-map)         | High     | Planned |
| [DNS Record Validator](#24-dns-record-validator)           | High     | Planned |
| [Remediation Guide](#25-remediation-guide)                 | High     | Planned |
| [IP Reputation Integration](#26-ip-reputation-integration) | Medium   | Planned |
| [Onboarding Wizard](#27-onboarding-wizard)                 | Medium   | Planned |
| [Webhook Support](#28-webhook-support)                     | Medium   | Planned |

---

### Phase 3: Enterprise Ready

**Theme**: Scale to teams, organizations, and enterprises

This phase introduces the collaboration and governance features needed for professional deployments. Parse DMARC becomes viable for organizations with compliance requirements and multiple stakeholders.

| Feature                                                      | Priority | Status  |
| ------------------------------------------------------------ | -------- | ------- |
| [Authentication System](#31-authentication-system)           | High     | Planned |
| [Multi-Organization Support](#32-multi-organization-support) | High     | Planned |
| [Role-Based Access Control](#33-role-based-access-control)   | Medium   | Planned |
| [Prometheus Metrics](#34-prometheus-metrics)                 | Medium   | Planned |
| [Scheduled Reports](#35-scheduled-reports)                   | Medium   | Planned |
| [Audit Logging](#36-audit-logging)                           | Medium   | Planned |
| [API Keys](#37-api-keys)                                     | Medium   | Planned |

---

### Phase 4: AI-Powered Security

**Theme**: Intelligence that actually helps

This phase leverages AI/ML to provide insights humans might miss and explanations that make complex concepts accessible. This is where Parse DMARC becomes truly differentiated.

| Feature                                                  | Priority | Status  |
| -------------------------------------------------------- | -------- | ------- |
| [AI Assistant](#41-ai-assistant)                         | Medium   | Planned |
| [Anomaly Detection](#42-anomaly-detection)               | Medium   | Planned |
| [Forensic Reports (RUF)](#43-forensic-reports)           | Medium   | Planned |
| [Natural Language Queries](#44-natural-language-queries) | Low      | Planned |
| [Predictive Analytics](#45-predictive-analytics)         | Low      | Planned |

---

## Feature Details

### Phase 1 Features

#### 1.1 Dark Mode

**Context**: Dark mode has evolved from a nice-to-have to a baseline expectation in 2025. Developers especially prefer dark interfaces, and DMARC monitoring often happens during incident response when eye strain matters.

**User Story**: As a user working late investigating an email issue, I want a dark interface option so that I can work comfortably without straining my eyes.

**Scope**:

- System preference detection (prefers-color-scheme)
- Manual toggle with persistence (localStorage)
- Smooth transition animations
- Consistent color palette across all components
- Proper contrast ratios for accessibility (WCAG AA)

**Technical Notes**:

- CSS custom properties for theming
- Single source of truth for color tokens
- Consider both light-on-dark and dark-on-light contrast needs

---

#### 1.2 DNS Record Generator

**Context**: The #1 barrier to DMARC adoption is the complexity of creating correct DNS records. Users must understand TXT record syntax, policy options, and reporting addresses. A generator eliminates this friction entirely.

**User Story**: As a new user, I want Parse DMARC to generate the exact DNS record I need so that I can copy-paste it into my DNS provider without learning DMARC syntax.

**Scope**:

- Interactive form for DMARC policy options:
  - Policy level (none/quarantine/reject)
  - Subdomain policy
  - Percentage
  - Alignment modes (relaxed/strict)
  - Reporting addresses (rua/ruf)
  - Reporting interval
- Real-time preview of generated record
- Copy-to-clipboard with visual feedback
- Provider-specific instructions (Cloudflare, Route53, GoDaddy, etc.)
- Validation of email addresses
- Explanation of each option in plain English

**Example Output**:

```
_dmarc.yourdomain.com TXT "v=DMARC1; p=reject; rua=mailto:dmarc@yourdomain.com; pct=100; adkim=r; aspf=r"
```

**Technical Notes**:

- Client-side generation (no server needed)
- Consider generating SPF records too as a bonus

---

#### 1.3 Compliance Health Score

**Context**: A single number (0-100) that represents overall email authentication health is instantly understandable by anyone—executives, developers, or security teams. It gamifies improvement and simplifies communication.

**User Story**: As a domain owner, I want a single score representing my email authentication health so that I can quickly understand my status and track improvement over time.

**Scope**:

- Score calculation algorithm considering:
  - DMARC pass rate (primary factor)
  - SPF alignment rate
  - DKIM alignment rate
  - Policy strictness (reject > quarantine > none)
  - Volume of authenticated traffic
- Visual representation (circular gauge, letter grade)
- Color coding (green/yellow/red zones)
- Trend indicator (improving/declining/stable)
- Breakdown of contributing factors
- Historical score tracking

**Scoring Framework**:

```
Base Score = (DMARC Pass Rate × 0.5) + (SPF Pass Rate × 0.25) + (DKIM Pass Rate × 0.25)
Policy Multiplier = { reject: 1.0, quarantine: 0.9, none: 0.7 }
Final Score = Base Score × Policy Multiplier × 100
```

---

#### 1.4 Data Export

**Context**: Users need to share DMARC data with stakeholders, import into other tools, or archive for compliance. Export functionality is table stakes for any monitoring tool.

**User Story**: As a security engineer, I want to export DMARC data in standard formats so that I can include it in reports, analyze in spreadsheets, or archive for compliance.

**Scope**:

- Export formats:
  - CSV (for spreadsheets)
  - JSON (for programmatic use)
  - PDF (for reports) - Phase 2
- Export scopes:
  - All reports
  - Filtered reports (date range, domain, etc.)
  - Single report detail
  - Statistics summary
- Column selection for CSV
- Date range filtering
- Immediate download (no email required)

**Technical Notes**:

- Client-side generation for CSV/JSON (faster, no server load)
- Consider streaming for large exports
- Include metadata (export date, filters applied)

---

#### 1.5 Keyboard Shortcuts

**Context**: Power users live by keyboard shortcuts. They signal a mature, thoughtfully-designed application and dramatically improve efficiency for frequent users.

**User Story**: As a power user, I want keyboard shortcuts for common actions so that I can navigate and operate the dashboard without reaching for my mouse.

**Scope**:

- Global shortcuts:
  - `r` - Refresh data
  - `/` or `Cmd+K` - Focus search/filter
  - `?` - Show keyboard shortcut help
  - `Esc` - Close modal/clear focus
- Navigation:
  - `j/k` - Navigate through report list
  - `Enter` - Open selected report
  - `g h` - Go to home/dashboard
- Help modal showing all shortcuts
- Shortcuts should not conflict with browser defaults
- Visible hints in UI (tooltips showing shortcuts)

**Technical Notes**:

- Use a lightweight library or native event handling
- Ensure shortcuts don't fire when typing in inputs
- Consider customization in future

---

#### 1.6 ASN/Organization Enrichment

**Context**: Raw IP addresses are meaningless to most users. Showing "AS15169 - Google LLC" instead of "172.217.164.110" provides instant context about whether a sender is legitimate.

**User Story**: As a user reviewing sending sources, I want to see organization names instead of just IP addresses so that I can quickly identify legitimate senders vs. suspicious ones.

**Scope**:

- Resolve IP → ASN → Organization name
- Display in:
  - Top sending sources table
  - Report detail view
  - Tooltips on IP addresses
- Caching to minimize lookups
- Graceful degradation when lookup fails
- Data sources:
  - MaxMind GeoLite2 ASN database (free)
  - Or ipinfo.io API (freemium)

**Example Display**:

```
172.217.164.110
Google LLC (AS15169)
Mountain View, US
```

---

#### 1.7 Compliance Badge

**Context**: Dynamic badges in README files are a powerful growth mechanism. Developers love showing off their infrastructure quality, and every badge is free advertising for Parse DMARC.

**User Story**: As a developer, I want a dynamic badge showing my domain's DMARC compliance so that I can display it in my README and demonstrate my commitment to email security.

**Scope**:

- Dynamic SVG badge endpoint: `/badge/{domain}.svg`
- Badge styles:
  - shields.io compatible format
  - Flat, flat-square, for-the-badge styles
- Shows compliance percentage and/or letter grade
- Color coded (green >90%, yellow >70%, red <70%)
- Caching for performance
- Rate limiting to prevent abuse
- Optional: require API key for private data

**Example**:

```markdown
![DMARC Compliance](https://your-instance.com/badge/example.com.svg)
```

**Security Consideration**: Consider whether compliance data should be public. May require opt-in or API key.

---

#### 1.8 Manual Report Upload

**Context**: Not everyone uses IMAP. Some users receive reports via other channels, want to test with sample data, or are evaluating Parse DMARC before configuring email access.

**User Story**: As a user without IMAP access, I want to upload DMARC XML files directly so that I can use Parse DMARC without configuring email integration.

**Scope**:

- Drag-and-drop upload zone
- Support formats:
  - Raw XML files
  - Gzip compressed (.xml.gz)
  - ZIP archives
- Multiple file upload
- Progress indication
- Duplicate detection (skip already-imported reports)
- Validation with helpful error messages
- Success confirmation with link to view report

**Technical Notes**:

- Reuse existing parser code
- Consider file size limits
- Add rate limiting for public instances

---

### Phase 2 Features

#### 2.1 Smart Alerting System

**Context**: A dashboard you have to check manually isn't monitoring—it's a chore. Real monitoring means being notified when something needs attention, without drowning in noise.

**User Story**: As a domain owner, I want to receive alerts when something important happens with my email authentication so that I can respond quickly to issues without constantly checking the dashboard.

**Scope**:

- Alert channels:
  - Email (via SMTP configuration)
  - Slack (webhook integration)
  - Discord (webhook integration)
  - Microsoft Teams (webhook integration)
  - Generic webhooks (for custom integrations)
- Alert types:
  - Compliance dropped below threshold
  - New sending source detected
  - Sudden volume spike/drop
  - Authentication failure rate increase
  - First report from new domain
- Configuration:
  - Per-channel settings
  - Alert frequency limits (no more than X per hour)
  - Quiet hours
  - Threshold customization
- Alert history and acknowledgment

**Example Alert**:

```
🚨 DMARC Alert: example.com

Compliance dropped from 98% to 72% in the last 24 hours.

Top failure source: 192.168.1.100 (Unknown - AS12345)
Failed messages: 1,247

View details: https://your-instance.com/reports?domain=example.com
```

---

#### 2.2 Historical Trend Charts

**Context**: Point-in-time metrics don't tell the story. Trends reveal whether you're improving, degrading, or stable—and help justify security investments to stakeholders.

**User Story**: As a security engineer, I want to see how my email authentication metrics have changed over time so that I can track improvement, detect gradual degradation, and report on progress.

**Scope**:

- Time-series charts for:
  - Compliance rate over time
  - Email volume trends
  - Pass/fail breakdown
  - Top senders over time
- Time range selection:
  - Last 7 days
  - Last 30 days
  - Last 90 days
  - Custom range
- Granularity options (daily, weekly)
- Multi-domain comparison
- Export charts as images
- Responsive design for all screen sizes

**Technical Notes**:

- Use a lightweight charting library (Chart.js, Apache ECharts)
- Aggregate data in database for performance
- Consider data retention implications

---

#### 2.3 Interactive GeoIP Map

**Context**: A world map showing where emails originate provides instant visual insight that tables cannot match. Geographic anomalies (emails from unexpected countries) become immediately obvious.

**User Story**: As a security analyst, I want to see a geographic visualization of email sources so that I can quickly identify suspicious activity from unexpected locations.

**Scope**:

- Interactive world map (Leaflet + OpenStreetMap)
- Features:
  - Markers/bubbles sized by volume
  - Color coding by authentication result
  - Click to drill down by country/region
  - Hover for details (country, volume, pass rate)
- Filtering:
  - By domain
  - By date range
  - By authentication result
- Heat map mode for high-volume data
- List view toggle for accessibility

**Technical Notes**:

- MaxMind GeoLite2 City database for IP geolocation
- Cluster markers at low zoom levels
- Consider WebGL for large datasets

---

#### 2.4 DNS Record Validator

**Context**: Many DMARC issues stem from misconfigured DNS records—syntax errors, missing components, conflicting policies. Proactive validation catches problems before they cause failures.

**User Story**: As a domain administrator, I want Parse DMARC to validate my DNS records and alert me to misconfigurations so that I can fix issues before they affect deliverability.

**Scope**:

- Validate records for:
  - DMARC (\_dmarc.domain.com)
  - SPF (domain.com TXT)
  - DKIM (selector.\_domainkey.domain.com)
- Validation checks:
  - Syntax correctness
  - Policy consistency
  - SPF record length (<255 chars, <10 DNS lookups)
  - DKIM key validity
  - Common mistakes detection
- Plain-English explanations of issues
- Severity levels (error, warning, info)
- Recommendations for fixes
- Scheduled periodic checks with alerting
- Manual on-demand validation

**Example Output**:

```
✅ DMARC: Valid
   Policy: reject | Alignment: relaxed

⚠️ SPF: Warning
   Record valid but approaching 10 DNS lookup limit (8/10)
   Recommendation: Consider flattening your SPF record

❌ DKIM: Error
   Selector 'google' not found
   Expected: google._domainkey.example.com
```

---

#### 2.5 Remediation Guide

**Context**: Knowing something failed isn't useful without knowing how to fix it. Context-aware remediation guidance transforms Parse DMARC from a diagnostic tool into a solution.

**User Story**: As a user viewing a failed authentication, I want specific guidance on how to fix the issue so that I don't have to research DMARC troubleshooting separately.

**Scope**:

- "Fix This" button on failed records
- Context-aware recommendations based on:
  - Failure type (SPF, DKIM, alignment)
  - Sender identification (known services vs. unknown)
  - Common patterns and known issues
- Guidance categories:
  - Legitimate sender misconfiguration
  - Third-party service setup
  - Potential spoofing/attack
  - DNS configuration issues
- Step-by-step instructions
- Links to relevant documentation
- Provider-specific guides (Google Workspace, Microsoft 365, SendGrid, etc.)

**Example**:

```
📧 Failure Analysis

Source: 198.51.100.50 (SendGrid, Inc.)
Issue: SPF alignment failed

This appears to be SendGrid sending on your behalf, but SPF
isn't properly configured.

Fix Steps:
1. Add SendGrid to your SPF record:
   include:sendgrid.net

2. Your new SPF record should be:
   v=spf1 include:_spf.google.com include:sendgrid.net ~all

3. Wait for DNS propagation (up to 48 hours)

📚 SendGrid DMARC Setup Guide →
```

---

#### 2.6 IP Reputation Integration

**Context**: An unknown IP sending email as your domain could be a new legitimate service—or an attacker. Reputation data provides immediate context that would otherwise require manual investigation.

**User Story**: As a security analyst, I want to see reputation information for sending IPs so that I can quickly distinguish between legitimate senders and potential threats.

**Scope**:

- Reputation data sources:
  - AbuseIPDB (abuse reports)
  - VirusTotal (malware/phishing associations)
  - Shodan (exposed services) - optional
  - Internal reputation (based on historical data)
- Display elements:
  - Reputation score/badge
  - Known associations (malware, spam, etc.)
  - Abuse report count
  - First/last seen dates
- Visual indicators:
  - 🟢 Clean/trusted
  - 🟡 Suspicious/limited data
  - 🔴 Known bad actor
- Caching to respect rate limits
- Graceful degradation without API keys

**Privacy Note**: Consider opt-in for external API calls that share IP data.

---

#### 2.7 Onboarding Wizard

**Context**: First impressions are permanent. A guided onboarding experience reduces time-to-value, prevents misconfiguration, and demonstrates that Parse DMARC is designed for humans.

**User Story**: As a new user, I want a guided setup experience so that I can get Parse DMARC working correctly without reading extensive documentation.

**Scope**:

- Step-by-step wizard:
  1. Welcome & overview
  2. DNS record setup (with generator)
  3. Email/IMAP configuration
  4. Connection test
  5. First report check
  6. Dashboard tour
- Progress indicator
- Skip option for experienced users
- Contextual help at each step
- Troubleshooting for common issues
- Completion celebration
- Option to revisit wizard

**Technical Notes**:

- Store completion state to not show repeatedly
- Consider video/animation elements
- Mobile-friendly design

---

#### 2.8 Webhook Support

**Context**: Webhooks enable Parse DMARC to be a trigger in larger automation workflows. Security teams can pipe events to SIEM, ticketing systems, or custom automation without polling.

**User Story**: As a DevOps engineer, I want to receive webhook notifications for DMARC events so that I can integrate with my existing automation and monitoring systems.

**Scope**:

- Webhook events:
  - `report.received` - New report processed
  - `compliance.changed` - Compliance threshold crossed
  - `source.new` - New sending source detected
  - `alert.triggered` - Any alert condition met
- Configuration:
  - Multiple webhook endpoints
  - Event filtering per endpoint
  - Custom headers (for auth)
  - Secret for payload signing (HMAC)
- Features:
  - Retry with exponential backoff
  - Delivery logging
  - Manual test/ping
  - Payload preview

**Payload Example**:

```json
{
  "event": "compliance.changed",
  "timestamp": "2025-11-28T10:30:00Z",
  "data": {
    "domain": "example.com",
    "previous_compliance": 98.5,
    "current_compliance": 72.3,
    "period": "24h",
    "top_failure_source": "192.168.1.100"
  }
}
```

---

### Phase 3 Features

#### 3.1 Authentication System

**Context**: Production deployments need access control. Without authentication, anyone who can reach the dashboard can view sensitive email security data.

**User Story**: As an administrator, I want to secure access to Parse DMARC so that only authorized users can view and manage email authentication data.

**Scope**:

- Authentication methods:
  - Local username/password
  - OIDC/OAuth2 (Google, GitHub, etc.)
  - SAML (enterprise SSO)
- Features:
  - Secure password hashing (argon2id)
  - Session management
  - Remember me option
  - Password reset flow
  - MFA support (TOTP)
- User management:
  - Invite users via email
  - Disable/enable accounts
  - Password policies

**Technical Notes**:

- Consider using an auth library (ZITADEL, Authelia, or built-in)
- JWT or session-based authentication
- Secure cookie handling

---

#### 3.2 Multi-Organization Support

**Context**: MSPs, agencies, and enterprises manage multiple domains across different clients or business units. A single Parse DMARC instance should support complete data isolation between organizations.

**User Story**: As an MSP, I want to manage multiple clients' DMARC data in a single Parse DMARC instance so that I don't need separate deployments for each client.

**Scope**:

- Organization model:
  - Create/edit/delete organizations
  - Assign domains to organizations
  - Separate IMAP credentials per org
  - Data isolation between orgs
- User-organization relationship:
  - Users belong to one or more orgs
  - Per-org roles
  - Organization switching in UI
- Per-organization settings:
  - Alert configurations
  - Retention policies
  - Branding (Phase 4)
- Global admin view across all orgs

**Database Considerations**:

- Add organization_id to reports/records tables
- Index optimization for multi-tenant queries
- Consider row-level security

---

#### 3.3 Role-Based Access Control

**Context**: Different users need different permissions. Admins configure; analysts investigate; executives view reports. RBAC enables secure delegation.

**User Story**: As an organization admin, I want to assign different permission levels to users so that I can control who can view, edit, or manage Parse DMARC.

**Scope**:

- Built-in roles:
  - **Owner**: Full access, can delete org
  - **Admin**: Manage users, settings, full data access
  - **Analyst**: View all data, create reports, manage alerts
  - **Viewer**: Read-only dashboard access
- Permissions:
  - View reports
  - Export data
  - Manage alerts
  - Manage users
  - Configure IMAP
  - Access API
- Custom roles (future)
- Role assignment UI
- Audit log of permission changes

---

#### 3.4 Prometheus Metrics

**Context**: DevOps teams expect metrics in Prometheus format for unified observability. Parse DMARC should integrate into existing monitoring stacks rather than requiring separate monitoring.

**User Story**: As a DevOps engineer, I want Parse DMARC to expose Prometheus metrics so that I can monitor it alongside my other infrastructure in Grafana.

**Scope**:

- Metrics endpoint: `/metrics`
- Metrics exposed:
  - `parse_dmarc_reports_total` (counter)
  - `parse_dmarc_messages_total{result="pass|fail"}` (counter)
  - `parse_dmarc_compliance_ratio{domain="..."}` (gauge)
  - `parse_dmarc_sources_unique` (gauge)
  - `parse_dmarc_fetch_duration_seconds` (histogram)
  - `parse_dmarc_fetch_errors_total` (counter)
- Labels for multi-domain filtering
- Standard Go runtime metrics
- Example Grafana dashboard JSON
- AlertManager rule examples

---

#### 3.5 Scheduled Reports

**Context**: Stakeholders who don't check dashboards still need visibility. Scheduled email reports keep everyone informed without requiring login access.

**User Story**: As a security manager, I want to receive scheduled summary reports so that I can keep stakeholders informed without them needing dashboard access.

**Scope**:

- Report types:
  - Daily summary
  - Weekly digest
  - Monthly executive report
- Delivery methods:
  - Email (SMTP)
  - Slack channel post
- Content:
  - Compliance trends
  - Notable events
  - Top issues requiring attention
  - Comparison to previous period
- Customization:
  - Recipients
  - Schedule (day/time)
  - Included domains
  - Report sections
- PDF attachment option

---

#### 3.6 Audit Logging

**Context**: Compliance frameworks (SOC 2, ISO 27001) require audit trails. Even without compliance requirements, knowing who did what and when is essential for security operations.

**User Story**: As a compliance officer, I want a complete audit trail of all actions in Parse DMARC so that I can demonstrate control effectiveness during audits.

**Scope**:

- Logged events:
  - User authentication (login, logout, failed attempts)
  - Configuration changes
  - User management actions
  - Data exports
  - Alert acknowledgments
  - API key creation/revocation
- Log details:
  - Timestamp
  - User/API key
  - Action type
  - Target resource
  - IP address
  - User agent
  - Before/after values (for changes)
- Log viewer in UI
- Export for external SIEM
- Retention configuration

---

#### 3.7 API Keys

**Context**: Programmatic access is essential for automation. API keys enable CI/CD integration, custom tooling, and third-party integrations without exposing user credentials.

**User Story**: As a developer, I want to generate API keys so that I can integrate Parse DMARC with my automation scripts and CI/CD pipelines.

**Scope**:

- API key management:
  - Create with name/description
  - Set expiration (optional)
  - Scope/permission limiting
  - Revoke immediately
  - View last used timestamp
- Security:
  - Keys shown only once at creation
  - Rate limiting per key
  - IP allowlisting (optional)
- Usage:
  - Bearer token authentication
  - Key prefix for identification (pd*live*...)
- Per-organization keys in multi-tenant mode

---

### Phase 4 Features

#### 4.1 AI Assistant

**Context**: DMARC is complex. Understanding why something failed, what it means, and what to do requires expertise. An AI assistant democratizes this expertise.

**User Story**: As a non-expert user, I want to ask questions about my DMARC data in plain English so that I can understand issues without becoming a DMARC expert.

**Scope**:

- Chat interface in dashboard
- Capabilities:
  - Explain what a specific failure means
  - Suggest fixes for common issues
  - Summarize recent activity
  - Answer DMARC/SPF/DKIM questions
  - Compare current state to best practices
- Context-aware (knows your data)
- Example queries:
  - "Why did this email fail?"
  - "Is this IP address legitimate?"
  - "What should I do about my SPF record?"
  - "Summarize last week's authentication issues"
- Privacy: Option to disable, clear data handling

**Technical Notes**:

- Integration with LLM API (OpenAI, Anthropic, local models)
- Careful prompt engineering for accuracy
- Rate limiting to control costs

---

#### 4.2 Anomaly Detection

**Context**: Manual threshold-based alerts miss subtle attacks. ML-based anomaly detection learns normal patterns and flags deviations humans might miss.

**User Story**: As a security engineer, I want Parse DMARC to automatically detect unusual patterns so that I can catch sophisticated attacks that wouldn't trigger simple threshold alerts.

**Scope**:

- Anomaly types:
  - Unusual volume (spike or drop)
  - New geographic source
  - Time-of-day anomalies
  - Sender pattern changes
  - Authentication result shifts
- Learning:
  - Baseline period (configurable)
  - Continuous learning
  - Per-domain models
- Alerting:
  - Anomaly score/confidence
  - Integration with alert system
  - Explanation of why flagged
- User feedback (mark as normal/expected)

---

#### 4.3 Forensic Reports

**Context**: Aggregate reports (RUA) show statistics; forensic reports (RUF) show individual message details. Complete DMARC support requires both.

**User Story**: As a security analyst investigating an incident, I want access to forensic reports so that I can see individual message details for failed authentications.

**Scope**:

- RUF report parsing
- Display:
  - Individual message headers
  - Authentication details
  - Message timestamps
  - Recipient information
- Privacy controls:
  - Redaction options
  - Retention limits
  - Access restrictions
- Correlation with aggregate reports
- Search by message ID, sender, etc.

**Note**: Forensic reports contain sensitive data. Implementation must prioritize privacy and security.

---

#### 4.4 Natural Language Queries

**Context**: Traditional filters require knowing the data model. Natural language lets users ask questions like they would ask a colleague.

**User Story**: As a user, I want to query my DMARC data using natural language so that I don't need to learn filter syntax.

**Scope**:

- Query examples:
  - "Show me failures from last week"
  - "Which IPs sent the most email in November?"
  - "Find reports from Google with high failure rates"
  - "Compare compliance this month vs last month"
- NLP to structured query translation
- Suggestions/autocomplete
- Save frequent queries
- Share queries with team

---

#### 4.5 Predictive Analytics

**Context**: Reactive security is always behind. Predicting compliance trends helps teams address issues before they become critical.

**User Story**: As a security manager, I want to see predictions of future compliance trends so that I can proactively address potential issues.

**Scope**:

- Predictions:
  - Compliance trajectory
  - Volume forecasting
  - Risk scoring for new sources
- Time horizon: 7-30 days
- Confidence intervals
- Factors influencing prediction
- Scenario modeling ("if we add this sender...")

---

## Future Exploration

These ideas are on our radar but not yet committed to the roadmap:

### Extended Protocol Support

- **MTA-STS monitoring**: Monitor Mail Transfer Agent Strict Transport Security policies
- **TLS-RPT parsing**: Parse TLS reporting alongside DMARC
- **BIMI validation**: Monitor Brand Indicators for Message Identification

### Advanced Integrations

- **Terraform provider**: Infrastructure-as-code for Parse DMARC configuration
- **Kubernetes operator**: CRD-based deployment and auto-scaling
- **SIEM connectors**: Native integration with Splunk, Elastic, Sentinel

### Community Features

- **Public threat intelligence feed**: Opt-in sharing of anonymized threat data
- **Community benchmarks**: Compare your metrics to anonymized industry averages
- **Plugin architecture**: Enable community-built extensions

### Innovative Concepts

- **Browser extension**: Check any domain's DMARC status with right-click
- **Slack bot**: Conversational interface for DMARC queries
- **Email deliverability score**: Holistic scoring beyond just DMARC

---

## Contributing

This roadmap is a living document. We welcome community input on prioritization and new ideas.

### How to Contribute

1. **Feature Requests**: Open a GitHub issue with the `enhancement` label
2. **Roadmap Discussion**: Comment on the roadmap discussion thread
3. **Implementation**: PRs welcome! Check issues labeled `help wanted`

### Prioritization Factors

We prioritize features based on:

- **Impact**: How many users benefit? How significant is the improvement?
- **Effort**: Engineering time and complexity
- **Strategic fit**: Alignment with project vision
- **Community demand**: GitHub issues, discussions, user feedback

---

## Changelog

| Date       | Change                                 |
| ---------- | -------------------------------------- |
| 2025-11-29 | DNS Record Generator feature completed |
| 2025-11-28 | Initial roadmap created                |

---

_This roadmap represents our current thinking and is subject to change based on community feedback, resource availability, and evolving priorities._
