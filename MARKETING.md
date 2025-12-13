# Parse DMARC - Open Source Marketing Strategy

> A comprehensive guide to promoting Parse DMARC across relevant platforms and communities.

## Target Audience

- **Primary**: DevOps engineers, SREs, System Administrators
- **Secondary**: Security professionals, Email administrators
- **Tertiary**: Self-hosters, Homelab enthusiasts, Small business IT

## Key Differentiators to Highlight

- Single 14MB binary/Docker image (vs Python + Elasticsearch + Kibana stack)
- Built-in Vue.js dashboard (no external visualization tools needed)
- SQLite storage (no JVM/Elasticsearch required)
- Zero dependencies, production-ready with Prometheus metrics
- Apache-2.0 licensed

---

## Tier 1: Highest-Intent Platforms (Submit First)

These platforms have audiences actively looking for self-hosted tools. **Start here.**

| Platform                  | Type            | URL                                                                                       | Submission      | Notes                                                                                                           |
| ------------------------- | --------------- | ----------------------------------------------------------------------------------------- | --------------- | --------------------------------------------------------------------------------------------------------------- |
| **awesome-selfhosted**    | GitHub List     | [github.com/awesome-selfhosted](https://github.com/awesome-selfhosted/awesome-selfhosted) | PR to repo      | **#1 priority** - Check CONTRIBUTING.md, submit under "Communication - Email - Mail Delivery Agents" or similar |
| **r/selfhosted**          | Reddit          | [reddit.com/r/selfhosted](https://reddit.com/r/selfhosted)                                | Post            | 400k+ members. Post as "New Release" with demo screenshot                                                       |
| **Hacker News (Show HN)** | Community       | [news.ycombinator.com](https://news.ycombinator.com/submit)                               | Direct submit   | Title: "Show HN: Parse DMARC - Single-binary DMARC report parser with built-in dashboard"                       |
| **awesome-go**            | GitHub List     | [github.com/avelino/awesome-go](https://github.com/avelino/awesome-go)                    | PR to repo      | Submit under "Email" or "Security" category                                                                     |
| **DevHunt**               | Launch Platform | [devhunt.org](https://devhunt.org)                                                        | Submit tool     | Free, 6-week queue (or $49 to skip). GitHub auth prevents bots                                                  |
| **Product Hunt**          | Launch Platform | [producthunt.com](https://producthunt.com)                                                | Schedule launch | Best on Tuesday-Thursday. Prep 250-char tagline + 5 screenshots                                                 |

### Submission Checklist for Tier 1

- [ ] awesome-selfhosted PR submitted
- [ ] r/selfhosted post published
- [ ] Show HN submitted
- [ ] awesome-go PR submitted
- [ ] DevHunt listing created
- [ ] Product Hunt launch scheduled

---

## Tier 2: High-Value Developer Communities

Active communities where DevOps/sysadmin content performs well.

| Platform          | Type            | URL                                                    | How to Submit | Notes                                                                      |
| ----------------- | --------------- | ------------------------------------------------------ | ------------- | -------------------------------------------------------------------------- |
| **Dev.to**        | Blog Platform   | [dev.to](https://dev.to)                               | Write article | Write "I built X to solve Y" post. Tags: #opensource #devops #security #go |
| **Lobsters**      | Link Aggregator | [lobste.rs](https://lobste.rs)                         | Invite-only   | Need invite. Technical focus. Mark as self-submission                      |
| **r/devops**      | Reddit          | [reddit.com/r/devops](https://reddit.com/r/devops)     | Post          | Share as tool announcement, focus on Prometheus metrics integration        |
| **r/sysadmin**    | Reddit          | [reddit.com/r/sysadmin](https://reddit.com/r/sysadmin) | Post          | Focus on email security monitoring angle                                   |
| **r/homelab**     | Reddit          | [reddit.com/r/homelab](https://reddit.com/r/homelab)   | Post          | Homelab-friendly: lightweight, Docker-ready                                |
| **r/golang**      | Reddit          | [reddit.com/r/golang](https://reddit.com/r/golang)     | Post          | Focus on Go implementation, architecture                                   |
| **Indie Hackers** | Community       | [indiehackers.com](https://www.indiehackers.com)       | Post on forum | Share the build story, open source journey                                 |
| **r/netsec**      | Reddit          | [reddit.com/r/netsec](https://reddit.com/r/netsec)     | Post          | Security angle: anti-phishing, email authentication                        |
| **r/docker**      | Reddit          | [reddit.com/r/docker](https://reddit.com/r/docker)     | Post          | Highlight 14MB image size, compose.yml                                     |

### Additional Reddit Communities

- r/homeserver
- r/kubernetes (for Helm chart/K8s deployment guides)
- r/Proxmox
- r/unRAID

---

## Tier 3: OSS Directories & Aggregators

Directories that drive SEO traffic and backlinks.

| Platform                          | URL                                              | Submission      | Cost | Notes                                                |
| --------------------------------- | ------------------------------------------------ | --------------- | ---- | ---------------------------------------------------- |
| **AlternativeTo**                 | [alternativeto.net](https://alternativeto.net)   | Add software    | Free | List as alternative to "ParseDMARC" (Python version) |
| **OpenAlternative**               | [openalternative.co](https://openalternative.co) | Submit          | Free | Open source alternatives directory                   |
| **LibHunt**                       | [libhunt.com](https://libhunt.com)               | Automatic       | Free | Tracks GitHub mentions. Add topics to repo           |
| **SaaSHub**                       | [saashub.com](https://www.saashub.com)           | Submit          | Free | DR 75, dofollow backlinks                            |
| **StackShare**                    | [stackshare.io](https://stackshare.io)           | Add tool        | Free | Good for "tech stack" discovery                      |
| **SourceForge**                   | [sourceforge.net](https://sourceforge.net)       | Create project  | Free | Legacy but still drives traffic                      |
| **OSS Software**                  | [osssoftware.org](https://osssoftware.org)       | Submit          | Free | Curated open source directory                        |
| **Slant**                         | [slant.co](https://slant.co)                     | Add option      | Free | "What's the best DMARC tool?" comparisons            |
| **Free Software Directory (FSF)** | [directory.fsf.org](https://directory.fsf.org)   | Wiki submission | Free | FSF-vetted, requires Apache-2.0 (compatible)         |

---

## Tier 4: Newsletters (Earned Media)

Getting featured in newsletters = high-quality traffic. Reach out to editors.

| Newsletter                | Focus            | URL                                                | How to Pitch                           |
| ------------------------- | ---------------- | -------------------------------------------------- | -------------------------------------- |
| **Golang Weekly**         | Go ecosystem     | [golangweekly.com](https://golangweekly.com)       | Contact Cooperpress, submit via site   |
| **DevOps Weekly**         | DevOps tools     | [devopsweekly.com](https://www.devopsweekly.com)   | Email Gareth Rushgrove                 |
| **Console.dev**           | Dev tools        | [console.dev](https://console.dev)                 | Submit via "Beta" section              |
| **Changelog Nightly**     | Trending repos   | [changelog.com](https://changelog.com)             | Automatic (tracks GitHub trending)     |
| **SRE Weekly**            | Site reliability | [sreweekly.com](https://sreweekly.com)             | Email editor with tool announcement    |
| **selfh.st**              | Self-hosted      | [selfh.st](https://selfh.st)                       | Contact for "This Week in Self-Hosted" |
| **tl;dr sec**             | Security         | [tldrsec.com](https://tldrsec.com)                 | Email Clint Gibler                     |
| **TLDR InfoSec**          | Security news    | [tldr.tech/infosec](https://tldr.tech/infosec)     | Submit via TLDR                        |
| **This Week in Security** | Cybersecurity    | [thisweekin.security](https://thisweekin.security) | Email Zack Whittaker                   |

### Newsletter Pitch Template

```
Subject: Open Source DMARC Parser - Single Binary, Built-in Dashboard

Hi [Name],

I built Parse DMARC, an open source tool that makes DMARC report analysis
accessible without the complexity of Elasticsearch stacks.

Key highlights:
- Single 14MB Docker image (vs Python + ES + Kibana)
- Built-in Vue.js dashboard
- Prometheus metrics for production monitoring
- Zero dependencies

GitHub: https://github.com/meysam81/parse-dmarc

Would this be a fit for [Newsletter Name]?

Best,
[Your Name]
```

---

## Tier 5: Security & Email-Specific Communities

Niche but highly relevant audiences.

| Platform                    | URL                                                                                                | Notes                                 |
| --------------------------- | -------------------------------------------------------------------------------------------------- | ------------------------------------- |
| **awesome-email-security**  | [github.com/0xAnalyst/awesome-email-security](https://github.com/0xAnalyst/awesome-email-security) | PR to add under DMARC tools           |
| **GitHub Topics**           | [github.com/topics/dmarc](https://github.com/topics/dmarc)                                         | Ensure repo has `dmarc` topic         |
| **GitHub Topics**           | [github.com/topics/email-security](https://github.com/topics/email-security)                       | Add `email-security` topic            |
| **DMARC.org**               | [dmarc.org](https://dmarc.org)                                                                     | Check for community resources section |
| **SpamAssassin Users List** | Mailing list                                                                                       | Email list mention                    |
| **Postfix Users**           | Mailing list                                                                                       | Relevant for mail server admins       |

### GitHub Topics to Add

Ensure your repo has these topics for discoverability:

- `dmarc`
- `email-security`
- `self-hosted`
- `spf`
- `dkim`
- `golang`
- `devops`
- `prometheus`
- `dashboard`
- `docker`

---

## Tier 6: Startup & Product Directories

Broader directories - less targeted but good for SEO.

| Platform           | URL                                            | Cost                       | Notes                                     |
| ------------------ | ---------------------------------------------- | -------------------------- | ----------------------------------------- |
| **BetaList**       | [betalist.com](https://betalist.com)           | Free (2mo wait) / $99 fast | Pre-launch focus                          |
| **BetaPage**       | [betapage.co](https://betapage.co)             | Free                       | Startup discovery                         |
| **StartupStash**   | [startupstash.com](https://startupstash.com)   | Free                       | Resources directory                       |
| **Launching Next** | [launchingnext.com](https://launchingnext.com) | Free                       | Daily curated list                        |
| **Capterra**       | [capterra.com](https://capterra.com)           | Free tier                  | Enterprise software directory             |
| **G2**             | [g2.com](https://g2.com)                       | Free tier                  | Software reviews (need users for reviews) |
| **GetApp**         | [getapp.com](https://getapp.com)               | Free tier                  | Gartner-owned directory                   |

---

## Tier 7: Social Media Strategy

Leverage existing social presence for announcements.

### Platforms & Content Strategy

| Platform         | Content Type       | Posting Strategy                                                    |
| ---------------- | ------------------ | ------------------------------------------------------------------- |
| **LinkedIn**     | Professional posts | Launch announcement, "Why I built this" story, technical deep-dives |
| **X (Twitter)**  | Threads & updates  | Launch thread, feature highlights, community engagement             |
| **Bluesky**      | Tech community     | Mirror X content, engage with self-hosting community                |
| **Mastodon**     | FOSS community     | Post on #opensource #selfhosted #foss hashtags                      |
| **Threads**      | Casual updates     | Repurpose LinkedIn/X content                                        |
| **Facebook**     | Groups             | Post in self-hosting, DevOps, homelab groups                        |
| **YouTube**      | Demo videos        | 2-5 min setup tutorial, dashboard walkthrough                       |
| **TikTok**       | Short clips        | "POV: You built a DMARC parser" style content                       |
| **BuyMeACoffee** | Project page       | [buymeacoffee.com](https://buymeacoffee.com) - Create project page  |
| **Patreon**      | Supporters         | Offer early access to features                                      |

### Sample Launch Tweet/Post

```
Just shipped Parse DMARC - an open source DMARC report parser
that doesn't require Elasticsearch.

- Single 14MB Docker image
- Built-in dashboard (no Kibana needed)
- Prometheus metrics included
- SQLite storage

Stop drowning in XML reports.

github.com/meysam81/parse-dmarc
```

### LinkedIn Post Structure

```
I built an open source tool to solve a problem that frustrated me for years.

DMARC reports are essential for email security, but:
- They arrive as compressed XML
- Existing tools need Elasticsearch + Kibana
- Setup takes hours, not minutes

So I built Parse DMARC:
- Single 14MB binary
- Vue.js dashboard included
- Prometheus metrics for prod
- Zero external dependencies

Open source. Apache-2.0 licensed.

Link in comments (GitHub)

#opensource #devops #emailsecurity #golang
```

---

## Tier 8: GitHub Optimization

Maximize discoverability within GitHub itself.

### Repository Checklist

- [ ] **Topics**: Add all relevant topics (see Tier 5)
- [ ] **Description**: Clear, keyword-rich one-liner
- [ ] **README badges**: License, Go Report Card, Docker pulls, GitHub stars
- [ ] **Demo screenshot**: High-quality dashboard image
- [ ] **Social preview image**: Custom Open Graph image (1280x640)
- [ ] **Releases**: Semantic versioning, detailed changelogs
- [ ] **Discussions**: Enable for community Q&A
- [ ] **Sponsor button**: Link to BuyMeACoffee/Patreon

### Trending Strategy

To hit GitHub Trending:

1. Coordinate social posts to drive stars in 24-48 hour window
2. Post on HN/Reddit at 9 AM EST (peak traffic)
3. Engage with every comment/issue promptly
4. Cross-post across all social channels simultaneously

---

## Tier 9: Discord Communities

Active developer communities on Discord.

| Server                       | Focus          | How to Join                                    |
| ---------------------------- | -------------- | ---------------------------------------------- |
| **r/selfhosted Discord**     | Self-hosting   | Link from subreddit sidebar                    |
| **r/homelab Discord**        | Homelab        | Link from subreddit sidebar                    |
| **Golang Discord**           | Go programming | [discord.gg/golang](https://discord.gg/golang) |
| **DevOps Discord**           | DevOps         | Search "DevOps" on Discord                     |
| **The Programmer's Hangout** | General dev    | Large community, share-your-projects channel   |

---

## Tier 10: Additional Awesome Lists

Submit PRs to relevant curated lists.

| List                   | URL                                                                                          | Category             |
| ---------------------- | -------------------------------------------------------------------------------------------- | -------------------- |
| **awesome-docker**     | [github.com/veggiemonk/awesome-docker](https://github.com/veggiemonk/awesome-docker)         | Monitoring/Security  |
| **awesome-sysadmin**   | [github.com/awesome-foss/awesome-sysadmin](https://github.com/awesome-foss/awesome-sysadmin) | Mail section         |
| **awesome-devops**     | Search GitHub                                                                                | DevOps tools         |
| **awesome-prometheus** | Search GitHub                                                                                | Prometheus exporters |
| **awesome-security**   | Search GitHub                                                                                | Email security tools |

---

## Execution Timeline

### Week 1: Foundation

- [ ] Optimize GitHub repo (topics, README, social preview)
- [ ] Submit to awesome-selfhosted
- [ ] Submit to awesome-go
- [ ] Post on r/selfhosted

### Week 2: Launch Push

- [ ] Schedule Product Hunt launch
- [ ] Submit to DevHunt
- [ ] Post Show HN
- [ ] Publish Dev.to article
- [ ] Social media launch posts (all platforms)

### Week 3: Directory Submissions

- [ ] Submit to all Tier 3 directories
- [ ] Submit to Tier 6 startup directories
- [ ] Submit to awesome-email-security

### Week 4: Outreach

- [ ] Email newsletter editors (Tier 4)
- [ ] Post in Discord communities
- [ ] Engage with Reddit comments
- [ ] Post on remaining subreddits (r/devops, r/sysadmin, r/netsec)

### Ongoing

- [ ] Respond to GitHub issues/discussions
- [ ] Share updates on social media
- [ ] Write follow-up Dev.to articles (tutorials, use cases)
- [ ] Create YouTube demo video
- [ ] Monitor and engage with community feedback

---

## Tracking & Metrics

### What to Track

- GitHub stars (primary metric)
- Docker pulls (ghcr.io)
- Website traffic (if applicable)
- Referral sources
- Newsletter features

### Tools

- GitHub Insights (stars, traffic, clones)
- Star-history.com for star growth visualization
- Social media analytics (native platforms)

---

## Resources

### Directory Aggregators

- [awesome-saas-directories](https://github.com/mahseema/awesome-saas-directories) - 100+ SaaS directories (some may accept open source/self-hosted tools; review for applicability)
- [Launchpedia](https://launchpedia.co) - Product Hunt alternatives list
- [Directory list by LinkDR](https://linkdr.com/blog/directories) - 99+ directories

### Inspiration

- How other self-hosted tools launched (Immich, Paperless-ngx, Vaultwarden)
- ParseDMARC (Python) marketing for comparison

---

_Last updated: December 2024_
_License: This document is part of the Parse DMARC project (Apache-2.0)_
