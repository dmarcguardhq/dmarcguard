# ChatGPT Brand Design Prompt for Parse DMARC

> Use this prompt with ChatGPT to generate a comprehensive brand design material stack.

---

## Context & Instructions

I'm building **Parse DMARC**, an open-source email authentication monitoring tool. I need you to act as my **product designer + brand strategist hybrid** and help me create a complete brand design system.

**This is executional** - I know what I'm building. I need cohesive brand identity, visual design system, and UI/UX refinement.

**Non-negotiable constraints:**
- Must work as a single-page dashboard application
- Open source project (Apache-2.0) - brand must feel welcoming to contributors
- Technical audience (developers, security engineers, DevOps) - no fluff, earned trust
- Small Docker image (14MB) is a key differentiator - simplicity is core to identity

---

## PART 1: PRODUCT CONTEXT

### What Parse DMARC Does
Parse DMARC is a self-hosted DMARC (Domain-based Message Authentication, Reporting & Conformance) report parser and monitoring dashboard. It:

1. **Fetches** DMARC aggregate reports from any IMAP mailbox (Gmail, Outlook, etc.)
2. **Parses** RFC 7489-compliant XML reports (including gzip/zip attachments)
3. **Stores** data in embedded SQLite (zero external dependencies)
4. **Displays** analytics in a Vue.js dashboard with real-time statistics
5. **Monitors** who's sending email on behalf of your domain

### The Problem We Solve
DMARC reports arrive as unreadable compressed XML attachments. Existing solutions (ParseDMARC Python) require Elasticsearch + Kibana + JVM - a massive stack for simple monitoring.

**Parse DMARC = single 14MB binary with built-in dashboard.** No external databases. No complex setup.

### Value Proposition
"Monitor who's sending email on behalf of your domain. Catch spoofing. Stop phishing. In a single 14MB binary."

### Target Users (Personas)

**Primary: Solo Developer / Startup Engineer**
- Runs multiple side projects with custom domains
- Wants set-and-forget DMARC monitoring
- Values simplicity over features
- Time-constrained, budget-conscious
- Success: Protected in under 5 minutes

**Secondary: Security Engineer**
- Protecting organization from spoofing/phishing
- Needs threat intelligence and alerting
- Values accuracy and actionable insights
- Integrates into existing security stack
- Success: Detects attacks in real-time

**Tertiary: MSP/Agency**
- Managing dozens/hundreds of client domains
- Needs multi-tenant capabilities
- Values unified visibility
- Requires white-label potential
- Success: Single pane of glass for 100+ domains

### Competitive Landscape
| Competitor | Parse DMARC Advantage |
|------------|----------------------|
| ParseDMARC (Python) | Single binary vs. Python + Elasticsearch + Kibana |
| OpenDMARC | Beautiful dashboard vs. raw CLI output |
| Commercial SaaS (Valimail, dmarcian) | Self-hosted, free, open source |
| Manual XML reading | Actually usable |

### Product Stage
**Early Mature / Growing**
- Version 1.3.10 (active development)
- Production deployments in the wild
- Prometheus metrics + Grafana dashboard included
- Multiple 1-click cloud deployments (Railway, Render, Koyeb, Zeabur)
- Homebrew distribution
- 14MB Docker image on Docker Hub

---

## PART 2: CURRENT VISUAL IDENTITY

### Existing Logo
Shield with checkmark - symbolizing email protection/verification
- Gradient fill: #667eea (indigo) → #764ba2 (purple)
- White checkmark inside shield
- Clean, modern, minimal

### Current Color Palette
```
Primary Gradient:    #667eea → #764ba2 (indigo to purple)
Success/Pass:        #4caf50 (green)
Error/Fail:          #f44336 (red)
Warning:             #fff3cd background, #856404 text
Neutral Text:        #333, #666, #999
Background:          Gradient header, white cards on subtle gradient body
```

### Current Typography
- System font stack: `-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Arial, sans-serif`
- Monospace for IPs/technical data: "Courier New", monospace
- Title: 2.5rem, font-weight 700
- Body: Default system sizing

### Current UI Patterns
- Glassmorphism effect on header/footer (rgba + backdrop-filter blur)
- White content cards with subtle box-shadows
- Rounded corners (12px cards, 6-8px buttons)
- Emoji as icons (🛡️ 📊 📧 ✅ 🌐 ⭐)
- Color-coded badges for compliance levels and policies

### Current Taglines
- "Email Authentication & Compliance Monitoring"
- "Monitor who's sending email on behalf of your domain. Catch spoofing. Stop phishing."
- "RFC 7489 Compliant DMARC Report Parser"

---

## PART 3: WHAT I NEED FROM YOU

### 1. Brand Foundations

**Naming Analysis**
- Is "Parse DMARC" the right name? Evaluate:
  - Memorability
  - SEO/searchability
  - Domain/trademark availability considerations
  - Alternatives if you see weaknesses

**Brand Voice & Messaging**
- Define the brand personality (3-5 adjectives with rationale)
- Writing style guide (tone, words to use, words to avoid)
- Elevator pitches: 1-sentence, 3-sentence, full paragraph versions
- Messaging hierarchy for different audiences

**Brand Positioning Statement**
- For [target user], Parse DMARC is the [category] that [key benefit] because [reason to believe]

### 2. Visual Identity System

**Logo Evolution**
- Critique the current shield+checkmark concept
- Propose refinements or alternatives (describe geometry, meaning, variants)
- Logo lockup variations needed:
  - Primary (icon + wordmark)
  - Icon only (for favicon, app icons)
  - Wordmark only
  - Dark mode variants
  - Monochrome versions

**Color System**
- Evaluate current palette - keep, refine, or replace?
- Define complete color system:
  - Primary, secondary, accent colors
  - Semantic colors (success, error, warning, info)
  - Neutral scale (text, backgrounds, borders)
  - Dark mode palette
- Provide hex values and usage guidelines
- Color accessibility check (WCAG AA compliance)

**Typography System**
- Recommend specific fonts (consider open source, web-safe options)
- Define type scale (headings, body, captions, code)
- Line heights, letter spacing
- Font weights and when to use each

**Spacing & Layout**
- Define spacing scale (base unit, multipliers)
- Grid system for dashboard layout
- Card patterns and composition rules
- Responsive breakpoints

### 3. Component Design System

**Core Components** (describe visual treatment for each):
- Buttons (primary, secondary, ghost, destructive)
- Cards (stat cards, report cards, detail cards)
- Badges/Pills (compliance levels, policy types, status indicators)
- Tables (sortable headers, row states, pagination)
- Forms (inputs, selects, toggles)
- Modal dialogs
- Navigation patterns
- Loading states
- Empty states
- Error states

**Data Visualization**
- Color encoding for pass/fail data
- Progress bars and gauges
- Compliance score visualization
- Source IP lists with reputation indicators
- Charts (if adding trend charts later)

**Iconography**
- Should we keep emojis or switch to icons?
- If icons: recommend icon set (style, source)
- Define usage rules

### 4. Key Screen Designs

Provide detailed design direction for these screens:

**A. Dashboard (Main View)**
- Statistics cards layout and hierarchy
- Top sending sources visualization
- Recent reports table
- Information architecture priorities
- What should draw the eye first?

**B. Report Detail Modal/Page**
- How to display authentication results clearly
- Per-record breakdown visualization
- Color coding for results
- Technical data presentation

**C. DNS Record Generator**
- Form UX best practices
- Live preview treatment
- Copy-to-clipboard interaction
- Provider-specific instructions display

**D. Empty State (No Reports Yet)**
- Onboarding guidance
- Call to action
- Reduce anxiety for new users

**E. Dark Mode**
- Complete dark theme treatment
- Maintain hierarchy and readability
- Handle gradients in dark mode

### 5. Brand Applications

**Marketing Assets**
- Open Graph image design (1200x630)
- GitHub social preview
- README header/banner
- Twitter/LinkedIn share cards
- Product Hunt launch graphics

**Documentation**
- Docs site visual treatment (if expanded)
- Code block styling
- Diagram style (architecture diagrams, flow charts)

**Community**
- GitHub repository presentation
- Badge designs for README
- Contributor-friendly visual language
- Sticker/swag concepts (optional fun)

### 6. Design Principles

Define 4-5 design principles that should guide all future design decisions. Format:
```
Principle Name
One-sentence description
What this means in practice
What this means we DON'T do
```

### 7. What NOT to Do

Create a "brand don'ts" list:
- Visual treatments to avoid
- Messaging anti-patterns
- UX patterns that contradict our values

---

## PART 4: EXISTING SCREENS FOR REFERENCE

### Dashboard Structure
```
┌─────────────────────────────────────────────────────────┐
│ HEADER (gradient background, glassmorphism)             │
│ 🛡️ DMARC Report Dashboard            [Refresh Button]   │
│ Email Authentication & Compliance Monitoring             │
├─────────────────────────────────────────────────────────┤
│ STAT CARDS (4-column grid)                              │
│ ┌───────┐ ┌───────┐ ┌───────┐ ┌───────┐                │
│ │📊     │ │📧     │ │✅     │ │🌐     │                │
│ │ 127   │ │45,892 │ │ 94.2% │ │  89   │                │
│ │Reports│ │Messages│ │Comply │ │Sources│                │
│ └───────┘ └───────┘ └───────┘ └───────┘                │
├─────────────────────────────────────────────────────────┤
│ TOP SENDING SOURCES                                     │
│ ┌─────────────────────────────────────────────────────┐│
│ │ 172.217.164.110  │ 12,456 msgs │ ████████░░ │ 89% ✓ ││
│ │ 209.85.220.41    │  8,234 msgs │ ███████░░░ │ 72% ✓ ││
│ │ 192.168.1.100    │  2,100 msgs │ ██░░░░░░░░ │ 23% ✗ ││
│ └─────────────────────────────────────────────────────┘│
├─────────────────────────────────────────────────────────┤
│ RECENT REPORTS (sortable table)                         │
│ Organization │ Domain │ Date │ Messages │ Comply │ Policy│
│ Google       │ ex.com │ 12/1 │ 1,234    │ 98.2%  │ reject│
│ Microsoft    │ ex.com │ 12/1 │   892    │ 95.1%  │ reject│
└─────────────────────────────────────────────────────────┘
```

### DNS Generator Structure
```
┌─────────────────────────────────────────────────────────┐
│ DMARC DNS Record Generator                              │
├──────────────────────┬──────────────────────────────────┤
│ POLICY OPTIONS       │ LIVE PREVIEW                     │
│                      │                                  │
│ Policy: [none ▼]     │ _dmarc.example.com TXT          │
│ Subdomain: [none ▼]  │ "v=DMARC1; p=reject;            │
│ Percentage: [100]    │  rua=mailto:dmarc@example.com;  │
│ RUA: [email input]   │  pct=100; adkim=r; aspf=r"      │
│ RUF: [email input]   │                                  │
│ DKIM Align: [relaxed]│              [📋 Copy]          │
│ SPF Align: [relaxed] │                                  │
├──────────────────────┴──────────────────────────────────┤
│ PROVIDER INSTRUCTIONS                                   │
│ [Cloudflare] [Route53] [GoDaddy] [Generic]             │
└─────────────────────────────────────────────────────────┘
```

---

## PART 5: TECHNICAL CONSTRAINTS

- Frontend: Vue.js 3 with Vite
- CSS: Scoped styles in Vue components (no external CSS framework)
- Icons: Currently emojis; could switch to SVG icon set
- Charts: None yet; considering Chart.js or Apache ECharts for Phase 2
- Dark mode: Needs CSS custom properties approach
- Must work in modern browsers (Chrome, Firefox, Safari, Edge)
- Should be responsive (desktop-first, but mobile-readable)
- Bundle size matters - avoid heavy dependencies

---

## PART 6: DELIVERABLES FORMAT

For each section, please provide:

1. **Rationale** - Why this choice? What problem does it solve?
2. **Specification** - Exact values, measurements, rules
3. **Examples** - Show me what it looks like in context
4. **Edge cases** - How does this work in unusual situations?
5. **Implementation notes** - Anything a developer needs to know

---

## PART 7: QUESTIONS TO ANSWER FIRST

Before diving into design, please address:

1. Is the current visual direction (purple gradient, shield icon, glassmorphism) on the right track, or should we pivot?
2. What's the biggest UX problem in the current design that we should fix first?
3. Is the emoji-as-icons approach working or should we professionalize with an icon set?
4. Should the brand feel more "security/enterprise" or "developer-friendly/approachable"?
5. Any red flags or missed opportunities in the current positioning?

---

## START HERE

Begin with the brand foundations (Part 3, Section 1), then work through systematically. Push back on any of my assumptions if you see problems. I want strategic clarity, not decoration.

Let's build something that makes email authentication actually feel accessible.
