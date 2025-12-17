# Parse DMARC - YouTube Video Script

## SEO Keyword Strategy

### Primary High-Intent Keywords (Target in Title/Description)

| Keyword                    | Search Intent  | Competition |
| -------------------------- | -------------- | ----------- |
| DMARC monitoring tool      | Solution-aware | Medium      |
| emails going to spam fix   | Problem-aware  | High        |
| self-hosted DMARC analyzer | Solution-aware | Low         |
| DMARC setup guide 2025     | How-to         | Medium      |
| email deliverability fix   | Problem-aware  | High        |
| read DMARC reports         | How-to         | Medium      |

### Adjacent Keywords (Natural inclusion in script)

- SPF DKIM DMARC setup
- Google Yahoo email requirements 2024/2025
- email authentication tutorial
- domain reputation fix
- bulk sender requirements
- open source email security
- homelab email monitoring
- DMARC dashboard
- email spoofing protection

### Long-tail Keywords (For description/tags)

- "why are my emails going to spam"
- "how to fix email deliverability"
- "free DMARC monitoring tool"
- "self-hosted alternative to EasyDMARC"
- "DMARC for small business"
- "Claude AI MCP integration"

---

## Video Metadata Recommendations

### Title Options (Pick One)

1. **"Stop Emails Going to Spam: Self-Hosted DMARC Monitoring (Free & Open Source)"** ← RECOMMENDED
2. "DMARC Monitoring in 5 Minutes: Free Tool That Actually Works"
3. "I Built a DMARC Dashboard So You Don't Have To (Self-Hosted)"
4. "Why Your Emails Land in Spam & How to Fix It (DMARC Tutorial 2025)"

### Description (Copy-paste ready)

```
Are your emails going to spam? With Google & Yahoo's new requirements and 87% of domains lacking proper DMARC protection, email deliverability is harder than ever.

In this video, I'll show you Parse DMARC - a free, open-source, self-hosted DMARC monitoring tool that:
✅ Fetches and parses DMARC reports automatically
✅ Beautiful dashboard to visualize your email security
✅ Works with Claude AI via MCP (Model Context Protocol)
✅ One Docker command to deploy
✅ No monthly fees or data limits

🔗 GitHub: https://github.com/meysam81/parse-dmarc
🐳 Docker: docker pull ghcr.io/meysam81/parse-dmarc

⏱️ Timestamps:
0:00 - Why your emails land in spam
1:45 - What is DMARC? (Quick explainer)
3:30 - The problem with existing tools
5:00 - Parse DMARC demo & features
7:30 - Installation (Docker one-liner)
9:00 - Reading your DMARC reports
10:30 - AI integration with Claude MCP
11:30 - Next steps & resources

📚 Resources mentioned:
- DMARC DNS Generator: Built into Parse DMARC
- Google's sender guidelines: https://support.google.com/mail/answer/81126
- Yahoo sender requirements: https://senders.yahooinc.com/

#DMARC #EmailSecurity #SelfHosted #OpenSource #EmailDeliverability #SPF #DKIM #Homelab #DevOps #EmailMarketing
```

### Tags (Copy-paste ready)

```
dmarc, dmarc monitoring, dmarc setup, email going to spam, fix email deliverability, spf dkim dmarc, email authentication, dmarc tutorial, self hosted, open source, email security, google email requirements, yahoo email requirements, dmarc reports, dmarc analyzer, homelab, devops, email marketing, bulk sender, domain reputation, email spoofing, phishing protection, parse dmarc, mcp claude ai
```

---

## Video Script

### HOOK (0:00 - 0:30)

**[Screen: Email in spam folder or "Message blocked" notification]**

```
Picture this: You've spent hours crafting the perfect email.
You hit send. And... nothing.

Your customer never sees it because it landed straight in spam.

Or worse - it got rejected entirely.

Here's the thing: 87% of domains don't have proper email authentication.
And with Google and Yahoo's new requirements, if you're not monitoring your
DMARC, you're basically flying blind.

Today I'm going to show you how to fix this - for free - with a self-hosted
tool that takes 30 seconds to deploy.
```

**[PATTERN INTERRUPT: Quick cut montage of dashboard, Docker command, reports]**

---

### PROBLEM AGITATION (0:30 - 1:45)

**[Screen: Statistics graphics, news headlines about email changes]**

```
Let me give you some numbers that should scare you:

77% of email deliverability problems come from bad domain reputation.

BEC attacks - that's business email compromise - have DOUBLED in the past year.
We're talking nearly 11 attacks per thousand mailboxes every month.

And here's the kicker: In February 2024, Google and Yahoo started requiring
DMARC for anyone sending more than 5,000 emails a day.

Microsoft just announced the same thing for May 2025.

If you're running a business, a newsletter, or even just a side project
that sends emails - this affects YOU.

The old approach of "set it and forget it" doesn't work anymore.

You need to MONITOR what's happening with your domain.

You need to KNOW when authentication fails.

And you need to do it without paying $50, $100, or even $300 a month
to some SaaS tool that's holding your data hostage.
```

---

### QUICK DMARC EXPLAINER (1:45 - 3:30)

**[Screen: Simple animated diagram or whiteboard explanation]**

```
Okay, quick primer if you're not familiar with DMARC.

Think of email authentication as a three-layer security system:

LAYER 1: SPF - Sender Policy Framework
This is basically a list of servers allowed to send email for your domain.
Like a bouncer with a VIP list.

LAYER 2: DKIM - DomainKeys Identified Mail
This cryptographically signs your emails.
Think of it as a wax seal proving the message wasn't tampered with.

LAYER 3: DMARC - Domain-based Message Authentication, Reporting & Conformance
This is the boss. It tells receiving servers:
"Hey, if an email fails SPF or DKIM, here's what to do with it."

And critically - DMARC sends you REPORTS.

Every day, Google, Yahoo, Microsoft - they all send you XML files
showing exactly what happened to emails claiming to be from your domain.

The problem? These reports look like THIS.

[Screen: Raw XML DMARC report - ugly, unreadable]

Yeah. Not exactly human-friendly.

Most people never read them. Which means they have no idea:
- Who's spoofing their domain
- Which legitimate services are failing authentication
- Whether their email is actually reaching inboxes

That's where Parse DMARC comes in.
```

---

### THE PROBLEM WITH EXISTING TOOLS (3:30 - 5:00)

**[Screen: Pricing pages of competitors, then Parse DMARC comparison]**

```
Now, there ARE tools out there to help with this.

EasyDMARC, PowerDMARC, dmarcian, Valimail...

They're fine. Some are even pretty good.

But here's what bugs me:

First - PRICING.
Most of these start at $15-50 per month for basic features.
Enterprise? You're looking at hundreds per month.
Per domain.

Second - DATA OWNERSHIP.
Your DMARC reports contain sensitive information about your email infrastructure.
Every IP that sends email for you. Every service. Every failure.
Do you really want that sitting on someone else's servers?

Third - LIMITS.
Free tiers are usually capped at like 10,000 messages or 2 domains.
Real businesses blow past that immediately.

And fourth - no AI integration.
We're in 2025. I want to ask my AI assistant:
"Hey, what's going on with my email authentication this week?"

That's why I started using Parse DMARC.

[Show GitHub page]

Open source. Self-hosted. Free forever.
Let me show you what it does.
```

---

### DEMO & FEATURES (5:00 - 7:30)

**[Screen: Live dashboard walkthrough]**

```
Alright, here's the Parse DMARC dashboard.

[Navigate through UI]

Right at the top - your compliance rate.
This tells you what percentage of emails are passing both SPF and DKIM.

You want this as close to 100% as possible.

Below that - total reports processed, total messages analyzed,
and a breakdown by source IP.

This is GOLD for debugging.

[Click into a report]

Here's an individual DMARC report from Google.
You can see:
- The reporting organization
- Date range covered
- Your published policy
- Each record showing source IP, message count, and pass/fail status

[Show top sources widget]

This widget shows your top sending sources.
Immediately you can see if there's an IP you don't recognize.
Spoofing attempt? Misconfigured service? Now you know.

[Show DNS Generator tool]

Oh, and it has this built-in DMARC DNS generator.
If you're setting up DMARC for the first time,
just fill in your email, pick your policy, and boom -
copy-paste the TXT record into your DNS.

No more googling "DMARC syntax" at 2am.

[Show recent reports list with pagination]

All your reports are stored locally in SQLite.
No cloud dependency. No API limits. Your data, your server.

And if you're a metrics person...

[Show /metrics endpoint]

Full Prometheus metrics endpoint.
Hook it up to Grafana, set up alerts.
Get notified the moment something goes wrong.
```

---

### INSTALLATION (7:30 - 9:00)

**[Screen: Terminal with Docker commands]**

```
Okay let's get this running. And I mean ACTUALLY running.

This isn't one of those "just follow these 47 steps" tutorials.

Ready?

[Type command]

docker run -d -p 8080:8080 \
  -e IMAP_HOST=imap.yourdomain.com \
  -e IMAP_USERNAME=dmarc@yourdomain.com \
  -e IMAP_PASSWORD=your-password \
  ghcr.io/meysam81/parse-dmarc:latest

That's it. That's the whole installation.

Now obviously you need a dedicated email address receiving your DMARC reports.

Quick tip: Create something like dmarc-reports@yourdomain.com
Then set your DMARC record's rua tag to point there.

[Show DMARC record example]

v=DMARC1; p=none; rua=mailto:dmarc-reports@yourdomain.com

The tool connects via IMAP, fetches the reports, parses them,
stores everything in a local SQLite database.

For production, you'll want to add a volume for persistence:

[Show docker-compose snippet]

There's a full docker-compose.yml in the repo
plus configs for Coolify, CapRover, even DigitalOcean droplets.

Whatever your deployment style, it's covered.
```

---

### READING YOUR REPORTS (9:00 - 10:30)

**[Screen: Dashboard with real report data - blur sensitive IPs if needed]**

```
So you've got it running. Reports are coming in.

Now what?

Let me show you how to actually READ this data.

[Point to compliance rate]

First - check your compliance rate.
Below 90%? You've got work to do.
100%? You're in good shape - but keep monitoring.

[Point to failing records]

Look for any records showing DKIM fail or SPF fail.

Common causes:
- A third-party service you forgot to add to SPF
- A marketing tool sending without proper DKIM
- Someone actually spoofing your domain

[Show IP lookup example]

Take that IP address, look it up.
Is it your server? Great, fix the config.
Is it AWS or Google Cloud? Probably a service you're using.
Is it some random hosting provider in another country?
Could be spoofing. Time to move your policy to quarantine or reject.

[Show trend over time]

Over time, you want to see:
- Compliance rate going UP
- Failed authentications going DOWN
- No surprise source IPs appearing

That's how you know your email authentication is solid.
```

---

### AI INTEGRATION - MCP (10:30 - 11:30)

**[Screen: Claude Desktop or similar with MCP integration]**

```
Now here's where it gets interesting.

Parse DMARC supports MCP - Model Context Protocol.

If you're using Claude Desktop or any MCP-compatible AI assistant,
you can query your DMARC data conversationally.

[Show MCP command or Claude interaction]

"What's my current DMARC compliance rate?"

"Show me any authentication failures from the past week"

"Which domains are sending the most email on my behalf?"

The AI can pull stats, analyze trends, even help you
troubleshoot specific failures.

This is running locally too - your data never leaves your server.

To enable it:

./parse-dmarc --mcp

For HTTP/SSE mode:

./parse-dmarc --mcp-http :8081

It even supports OAuth if you need authentication.

This is honestly the feature that got me excited about this tool.
Combining security monitoring with AI assistance
is exactly the kind of workflow I want in 2025.
```

---

### CALL TO ACTION & WRAP-UP (11:30 - 12:30)

**[Screen: GitHub page, then face cam or logo]**

```
Alright, let's wrap this up.

If you're sending emails - for your business, your side project,
your newsletter, whatever - you NEED DMARC monitoring.

Not setting up DMARC is like leaving your front door unlocked
and hoping nobody notices.

Parse DMARC gives you:
- A beautiful dashboard to visualize your email security
- Automatic report fetching and parsing
- Prometheus metrics for alerting
- AI integration via MCP
- And it's completely free and open source

Link to the GitHub is in the description.
Star it if you find it useful.

If you have questions, drop them in the comments.
I'll do my best to answer.

And if you want to see more self-hosted tools
and email security content, hit subscribe.

Your emails deserve to land in inboxes, not spam folders.

Go fix your DMARC.

See you in the next one.
```

---

## B-Roll Shot List

| Timestamp   | Shot Description                                |
| ----------- | ----------------------------------------------- |
| 0:00-0:15   | Email landing in spam folder (screen recording) |
| 0:15-0:30   | Quick cuts: Dashboard → Docker → Terminal       |
| 0:45-1:00   | Statistics graphics (animated text overlays)    |
| 1:45-3:30   | Animated diagram: SPF → DKIM → DMARC flow       |
| 3:00-3:15   | Raw XML DMARC report (ugly code view)           |
| 3:30-5:00   | Competitor pricing pages (blur logos if needed) |
| 5:00-7:30   | Full dashboard walkthrough (live recording)     |
| 7:30-9:00   | Terminal: Docker commands being typed           |
| 9:00-10:30  | Dashboard: Report analysis                      |
| 10:30-11:30 | MCP/AI demo (Claude Desktop or similar)         |
| 11:30-12:30 | GitHub repo, Star button, Subscribe animation   |

---

## Thumbnail Concepts

**Option A (Problem-focused):**

- Split screen: Left = sad face + "SPAM" red stamp, Right = happy face + inbox
- Text: "FIX DMARC" in bold
- Small: "FREE TOOL" badge

**Option B (Tool-focused):**

- Dashboard screenshot with blur effect
- Face cam reaction (surprised/excited)
- Text: "This FREE Tool..."

**Option C (Stats hook):**

- Large "87%" text
- Subtext: "of domains are VULNERABLE"
- Arrow pointing to dashboard preview

---

## Production Notes

1. **Video Length Target:** 10-12 minutes (sweet spot for retention + ad revenue)
2. **Hook:** First 30 seconds determine retention - make it punchy
3. **Pattern Interrupts:** Every 2-3 minutes, change visual format
4. **End Screen:** Last 20 seconds for subscribe/video cards
5. **Chapters:** Use timestamps in description for better SEO
6. **Cards:** Link to related videos at 3:00 (DMARC basics) and 9:00 (Docker tutorial)

---

## Post-Publish Checklist

- [ ] Pin comment with GitHub link and timestamps
- [ ] Reply to first 10 comments within 24 hours
- [ ] Share to relevant subreddits: r/selfhosted, r/homelab, r/devops, r/sysadmin
- [ ] Cross-post to Twitter/X with key stat hook
- [ ] Submit to Hacker News (title: "Show HN: Open source DMARC monitoring dashboard")
- [ ] Add to relevant Discord communities (homelab, DevOps)

---

## Sources & References

Research sources used for keyword analysis and content:

- [Email Vendor Selection - Best DMARC Monitoring Tools](https://www.emailvendorselection.com/best-dmarc-monitoring-tools/)
- [PowerDMARC - How to Read DMARC Reports](https://powerdmarc.com/how-to-read-dmarc-reports/)
- [MailReach - Fix Email Deliverability Issues](https://www.mailreach.co/blog/fix-email-deliverability-issues)
- [EasyDMARC - DMARC Step-by-Step Guide](https://easydmarc.com/blog/dmarc-step-by-step-guide/)
- [Mimecast - How to Read DMARC Reports](https://www.mimecast.com/blog/how-to-read-dmarc-reports/)
- [Cloudflare - DMARC DKIM SPF Explained](https://www.cloudflare.com/learning/email-security/dmarc-dkim-spf/)
- [Google Workspace - Set up DMARC](https://support.google.com/a/answer/2466580)
- [Mailtrap - Email Domain Reputation](https://mailtrap.io/blog/email-domain-reputation/)
- [Stalwart Mail Server](https://stalw.art/)
- [Mox Mail Server](https://github.com/mjl-/mox)
