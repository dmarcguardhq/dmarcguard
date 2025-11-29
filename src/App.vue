<template>
  <div class="app">
    <header class="header">
      <div class="container">
        <div class="header-top">
          <div>
            <h1 class="title">
              <span class="icon">🛡️</span>
              DMARC Report Dashboard
            </h1>
            <p class="subtitle">Email Authentication & Compliance Monitoring</p>
          </div>
          <div class="header-actions">
            <button
              class="theme-toggle"
              @click="toggleTheme"
              :title="theme === 'dark' ? 'Switch to light mode' : 'Switch to dark mode'"
              :aria-label="theme === 'dark' ? 'Switch to light mode' : 'Switch to dark mode'"
              :aria-pressed="theme === 'dark'"
            >
              <span v-if="theme === 'dark'" class="theme-icon" aria-hidden="true">
                <!-- Sun SVG icon -->
                <svg width="20" height="20" viewBox="0 0 20 20" fill="none" aria-hidden="true" focusable="false">
                  <circle cx="10" cy="10" r="5" fill="currentColor"/>
                  <g stroke="currentColor" stroke-width="2">
                    <line x1="10" y1="1" x2="10" y2="3"/>
                    <line x1="10" y1="17" x2="10" y2="19"/>
                    <line x1="1" y1="10" x2="3" y2="10"/>
                    <line x1="17" y1="10" x2="19" y2="10"/>
                    <line x1="4.22" y1="4.22" x2="5.64" y2="5.64"/>
                    <line x1="14.36" y1="14.36" x2="15.78" y2="15.78"/>
                    <line x1="4.22" y1="15.78" x2="5.64" y2="14.36"/>
                    <line x1="14.36" y1="5.64" x2="15.78" y2="4.22"/>
                  </g>
                </svg>
                <span class="visually-hidden">Light mode</span>
              </span>
              <span v-else class="theme-icon" aria-hidden="true">
                <!-- Moon SVG icon -->
                <svg width="20" height="20" viewBox="0 0 20 20" fill="none" aria-hidden="true" focusable="false">
                  <path
                    d="M15.5 13.5A7 7 0 0 1 6.5 4.5a7 7 0 1 0 9 9z"
                    fill="currentColor"
                  />
                </svg>
                <span class="visually-hidden">Dark mode</span>
              </span>
            </button>
            <button
              class="refresh-button"
              @click="refreshData"
              :disabled="loading"
            >
              <span class="refresh-icon" :class="{ spinning: loading }">🔄</span>
              <span>{{ loading ? "Refreshing..." : "Refresh" }}</span>
            </button>
          </div>
        </div>
      </div>
    </header>

    <!-- Star Banner -->
    <div v-if="starBannerVisible" class="star-banner">
      <div class="container">
        <div class="star-banner-content">
          <div class="star-banner-text">
            <span class="star-banner-icon">⭐</span>
            <span class="star-banner-message">
              <strong>Self-hosting Parse DMARC?</strong> Support open-source
              email security! Give us a star on GitHub and help others discover
              this tool.
            </span>
          </div>
          <div class="star-banner-actions">
            <a
              href="https://github.com/meysam81/parse-dmarc"
              target="_blank"
              rel="noopener noreferrer"
              class="star-button"
            >
              ⭐ Star on GitHub
            </a>
            <button
              class="dismiss-button"
              @click="dismissStarBanner"
              title="Dismiss star banner"
              aria-label="Dismiss star banner"
            >
              ×
            </button>
          </div>
        </div>
      </div>
    </div>

    <main class="main">
      <div class="container">
        <!-- Statistics Cards -->
        <div class="stats-grid" v-if="statistics">
          <div class="stat-card">
            <div class="stat-icon">📊</div>
            <div class="stat-content">
              <div class="stat-value">{{ statistics.total_reports }}</div>
              <div class="stat-label">Total Reports</div>
            </div>
          </div>

          <div class="stat-card">
            <div class="stat-icon">📧</div>
            <div class="stat-content">
              <div class="stat-value">
                {{ formatNumber(statistics.total_messages) }}
              </div>
              <div class="stat-label">Total Messages</div>
            </div>
          </div>

          <div class="stat-card">
            <div class="stat-icon">✅</div>
            <div class="stat-content">
              <div class="stat-value">
                {{ statistics.compliance_rate.toFixed(1) }}%
              </div>
              <div class="stat-label">Compliance Rate</div>
            </div>
          </div>

          <div class="stat-card">
            <div class="stat-icon">🌐</div>
            <div class="stat-content">
              <div class="stat-value">{{ statistics.unique_source_ips }}</div>
              <div class="stat-label">Unique Sources</div>
            </div>
          </div>
        </div>

        <!-- Top Sources -->
        <div class="section">
          <h2 class="section-title">Top Sending Sources</h2>
          <div class="card">
            <div class="source-list" v-if="topSources?.length > 0">
              <div
                v-for="source in topSources"
                :key="source.source_ip"
                class="source-item"
              >
                <div class="source-ip">{{ source.source_ip }}</div>
                <div class="source-stats">
                  <div class="source-count">
                    {{ formatNumber(source.count) }} messages
                  </div>
                  <div class="source-bar">
                    <div
                      class="source-bar-pass"
                      :style="{ width: getPassPercentage(source) + '%' }"
                    ></div>
                    <div
                      class="source-bar-fail"
                      :style="{ width: getFailPercentage(source) + '%' }"
                    ></div>
                  </div>
                  <div class="source-legend">
                    <span class="legend-pass">{{ source.pass }} pass</span>
                    <span class="legend-fail">{{ source.fail }} fail</span>
                  </div>
                </div>
              </div>
            </div>
            <div v-else class="empty-state">
              No data available yet. Reports will appear here once fetched.
            </div>
          </div>
        </div>

        <!-- Recent Reports -->
        <div class="section">
          <h2 class="section-title">Recent Reports</h2>
          <div class="card">
            <div class="table-container" v-if="reports?.length > 0">
              <table class="report-table">
                <thead>
                  <tr>
                    <th @click="sortBy('org_name')" class="sortable">
                      Organization
                      <span class="sort-indicator">{{
                        getSortIndicator("org_name")
                      }}</span>
                    </th>
                    <th @click="sortBy('domain')" class="sortable">
                      Domain
                      <span class="sort-indicator">{{
                        getSortIndicator("domain")
                      }}</span>
                    </th>
                    <th @click="sortBy('date_begin')" class="sortable">
                      Date Range
                      <span class="sort-indicator">{{
                        getSortIndicator("date_begin")
                      }}</span>
                    </th>
                    <th @click="sortBy('total_messages')" class="sortable">
                      Messages
                      <span class="sort-indicator">{{
                        getSortIndicator("total_messages")
                      }}</span>
                    </th>
                    <th @click="sortBy('compliance_rate')" class="sortable">
                      Compliance
                      <span class="sort-indicator">{{
                        getSortIndicator("compliance_rate")
                      }}</span>
                    </th>
                    <th @click="sortBy('policy_p')" class="sortable">
                      Policy
                      <span class="sort-indicator">{{
                        getSortIndicator("policy_p")
                      }}</span>
                    </th>
                  </tr>
                </thead>
                <tbody>
                  <tr
                    v-for="report in sortedReports"
                    :key="report.id"
                    class="report-row"
                    @click="viewReport(report)"
                  >
                    <td>{{ report.org_name }}</td>
                    <td>
                      <code>{{ report.domain }}</code>
                    </td>
                    <td class="date-cell">
                      {{ formatDate(report.date_begin) }}
                    </td>
                    <td>{{ formatNumber(report.total_messages) }}</td>
                    <td>
                      <span
                        class="compliance-badge"
                        :class="getComplianceClass(report.compliance_rate)"
                      >
                        {{ report.compliance_rate.toFixed(1) }}%
                      </span>
                    </td>
                    <td>
                      <span class="policy-badge">{{ report.policy_p }}</span>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
            <div v-else class="empty-state">
              No reports available yet. Check your IMAP configuration and run
              the fetch process.
            </div>
          </div>
        </div>

        <!-- Report Detail Modal -->
        <div v-if="selectedReport" class="modal" @click="closeModal">
          <div class="modal-content" @click.stop>
            <div class="modal-header">
              <h3>Report Details</h3>
              <button class="modal-close" @click="closeModal">×</button>
            </div>
            <div class="modal-body">
              <div class="detail-grid">
                <div class="detail-item">
                  <strong>Organization:</strong>
                  {{ selectedReport.ReportMetadata?.OrgName || "N/A" }}
                </div>
                <div class="detail-item">
                  <strong>Domain:</strong>
                  {{ selectedReport.PolicyPublished?.Domain || "N/A" }}
                </div>
                <div class="detail-item">
                  <strong>Report ID:</strong>
                  {{ selectedReport.ReportMetadata?.ReportID || "N/A" }}
                </div>
                <div class="detail-item">
                  <strong>Policy:</strong>
                  {{ selectedReport.PolicyPublished?.P || "N/A" }}
                </div>
              </div>

              <h4 class="detail-subtitle">
                Records ({{ selectedReport?.Records?.length || 0 }})
              </h4>
              <div class="records-list">
                <div
                  v-for="(record, idx) in selectedReport.Records || []"
                  :key="idx"
                  class="record-item"
                >
                  <div class="record-header">
                    <span class="record-ip">{{
                      record.Row?.SourceIP || "N/A"
                    }}</span>
                    <span class="record-count"
                      >{{ record.Row?.Count || 0 }} messages</span
                    >
                  </div>
                  <div class="record-details">
                    <span
                      :class="
                        'result-badge ' +
                        (record.Row?.PolicyEvaluated?.DKIM === 'pass'
                          ? 'pass'
                          : 'fail')
                      "
                    >
                      DKIM: {{ record.Row?.PolicyEvaluated?.DKIM || "unknown" }}
                    </span>
                    <span
                      :class="
                        'result-badge ' +
                        (record.Row?.PolicyEvaluated?.SPF === 'pass'
                          ? 'pass'
                          : 'fail')
                      "
                    >
                      SPF: {{ record.Row?.PolicyEvaluated?.SPF || "unknown" }}
                    </span>
                    <span class="result-badge">
                      Disposition:
                      {{
                        record.Row?.PolicyEvaluated?.Disposition || "unknown"
                      }}
                    </span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>

    <footer class="footer">
      <div class="container">
        <div class="footer-content">
          <div class="footer-section">
            <h4>Parse DMARC</h4>
            <p>RFC 7489 Compliant DMARC Report Parser</p>
          </div>
          <div class="footer-section">
            <h4>Links</h4>
            <div class="footer-links">
              <a
                href="https://github.com/meysam81/parse-dmarc"
                target="_blank"
                rel="noopener"
                >GitHub</a
              >
              <a
                href="https://github.com/meysam81/parse-dmarc/issues"
                target="_blank"
                rel="noopener"
                >Issues</a
              >
              <a
                href="https://github.com/meysam81/parse-dmarc/blob/main/README.md"
                target="_blank"
                rel="noopener"
                >Documentation</a
              >
            </div>
          </div>
          <div class="footer-section">
            <h4>Resources</h4>
            <div class="footer-links">
              <a
                href="https://datatracker.ietf.org/doc/html/rfc7489"
                target="_blank"
                rel="noopener"
                >RFC 7489</a
              >
              <a href="https://dmarc.org/" target="_blank" rel="noopener"
                >DMARC.org</a
              >
              <a
                href="https://github.com/meysam81/parse-dmarc/blob/main/LICENSE"
                target="_blank"
                rel="noopener"
                >Apache-2.0 License</a
              >
            </div>
          </div>
        </div>
        <div class="footer-bottom">
          <p>
            Built with
            <a href="https://vuejs.org/" target="_blank" rel="noopener"
              >Vue 3</a
            >
            +
            <a href="https://golang.org/" target="_blank" rel="noopener">Go</a>
          </p>
          <p class="opensource-message">
            Free and Open Source Software • Contributions Welcome • Made with ❤️
            for the Community
          </p>
        </div>
      </div>
    </footer>
  </div>
</template>

<script>
import { computed, onMounted, ref, watchEffect } from "vue";

export default {
  name: "App",
  setup() {
    var statistics = ref(null);
    var topSources = ref([]);
    var reports = ref([]);
    var selectedReport = ref(null);
    var loading = ref(true);
    var sortColumn = ref(null);
    var sortDirection = ref("asc");
    var starBannerVisible = ref(true);
    var theme = ref("light");

    function getSystemTheme() {
      if (
        typeof window !== "undefined" &&
        typeof window.matchMedia === "function" &&
        window.matchMedia("(prefers-color-scheme: dark)").matches
      ) {
        return "dark";
      }
      return "light";
    }

    function initTheme() {
      var savedTheme = localStorage.getItem("parse-dmarc-theme");
      if (savedTheme && (savedTheme === "light" || savedTheme === "dark")) {
        theme.value = savedTheme;
      } else {
        theme.value = getSystemTheme();
      }
    }

    function toggleTheme() {
      theme.value = theme.value === "dark" ? "light" : "dark";
      localStorage.setItem("parse-dmarc-theme", theme.value);
    }

    // Apply theme to document
    watchEffect(function () {
      document.documentElement.setAttribute("data-theme", theme.value);
    });

    function fetchStatistics() {
      return fetch("/api/statistics")
        .then(function (response) {
          return response.json();
        })
        .then(function (data) {
          statistics.value = data;
        })
        .catch(function (error) {
          console.error("Failed to fetch statistics:", error);
        });
    }

    function fetchTopSources() {
      return fetch("/api/top-sources?limit=10")
        .then(function (response) {
          return response.json();
        })
        .then(function (data) {
          topSources.value = data;
        })
        .catch(function (error) {
          console.error("Failed to fetch top sources:", error);
        });
    }

    function fetchReports() {
      return fetch("/api/reports?limit=20")
        .then(function (response) {
          return response.json();
        })
        .then(function (data) {
          reports.value = data;
        })
        .catch(function (error) {
          console.error("Failed to fetch reports:", error);
        });
    }

    function viewReport(report) {
      fetch(`/api/reports/${report.id}`)
        .then(function (response) {
          return response.json();
        })
        .then(function (data) {
          selectedReport.value = data;
        })
        .catch(function (error) {
          console.error("Failed to fetch report details:", error);
          selectedReport.value = null;
        });
    }

    function closeModal() {
      selectedReport.value = null;
    }

    function formatNumber(num) {
      return new Intl.NumberFormat().format(num);
    }

    function formatDate(timestamp) {
      return new Date(timestamp * 1000).toLocaleDateString();
    }

    function getPassPercentage(source) {
      var total = source.count;
      return total > 0 ? (source.pass / total) * 100 : 0;
    }

    function getFailPercentage(source) {
      var total = source.count;
      return total > 0 ? (source.fail / total) * 100 : 0;
    }

    function getComplianceClass(rate) {
      if (rate >= 95) return "high";
      if (rate >= 70) return "medium";
      return "low";
    }

    function loadData() {
      loading.value = true;
      return Promise.all([
        fetchStatistics(),
        fetchTopSources(),
        fetchReports(),
      ]).then(function () {
        loading.value = false;
      });
    }

    function refreshData() {
      loadData();
    }

    function sortBy(column) {
      if (sortColumn.value === column) {
        if (sortDirection.value === "asc") {
          sortDirection.value = "desc";
        } else {
          sortColumn.value = null;
          sortDirection.value = "asc";
        }
      } else {
        sortColumn.value = column;
        sortDirection.value = "asc";
      }
    }

    function getSortIndicator(column) {
      if (sortColumn.value !== column) {
        return "";
      }
      return sortDirection.value === "asc" ? "↑" : "↓";
    }

    var sortedReports = computed(function () {
      if (!sortColumn.value || !reports.value) {
        return reports.value;
      }

      var sorted = [...reports.value].sort(function (a, b) {
        var aVal = a[sortColumn.value];
        var bVal = b[sortColumn.value];

        if (typeof aVal === "string") {
          aVal = aVal.toLowerCase();
          bVal = bVal.toLowerCase();
        }

        if (aVal < bVal) {
          return sortDirection.value === "asc" ? -1 : 1;
        }
        if (aVal > bVal) {
          return sortDirection.value === "asc" ? 1 : -1;
        }
        return 0;
      });

      return sorted;
    });

    function dismissStarBanner() {
      starBannerVisible.value = false;
      var dismissalData = {
        dismissed: true,
        timestamp: Date.now(),
      };
      localStorage.setItem(
        "starBannerDismissed",
        JSON.stringify(dismissalData),
      );
    }

    function checkStarBannerDismissal() {
      var dismissalData = localStorage.getItem("starBannerDismissed");
      if (dismissalData) {
        try {
          var data = JSON.parse(dismissalData);
          var thirtyDaysInMs = 30 * 24 * 60 * 60 * 1000;
          if (
            data &&
            typeof data.timestamp === "number" &&
            !isNaN(data.timestamp)
          ) {
            var timeSinceDismissal = Date.now() - data.timestamp;
            if (timeSinceDismissal < thirtyDaysInMs) {
              starBannerVisible.value = false;
            }
          } else {
            // Clear old dismissal data after 30 days and show the banner again
            localStorage.removeItem("starBannerDismissed");
            starBannerVisible.value = true;
          }
        } catch (e) {
          console.error("Failed to parse star banner dismissal data:", e);
        }
      }
    }

    onMounted(function () {
      initTheme();
      checkStarBannerDismissal();
      loadData();
      setInterval(loadData, 5 * 60 * 1000);
    });

    return {
      statistics,
      topSources,
      reports,
      selectedReport,
      loading,
      sortColumn,
      sortDirection,
      sortedReports,
      starBannerVisible,
      theme,
      viewReport,
      closeModal,
      formatNumber,
      formatDate,
      getPassPercentage,
      getFailPercentage,
      getComplianceClass,
      sortBy,
      getSortIndicator,
      refreshData,
      dismissStarBanner,
      toggleTheme,
    };
  },
};
</script>

<style scoped>
/* Theme Color Tokens */
:root {
  /* Background colors */
  --color-bg-primary: #ffffff;
  --color-bg-secondary: #f8f9fa;
  --color-bg-tertiary: #e9ecef;
  --color-bg-header: rgba(255, 255, 255, 0.1);
  --color-bg-overlay: rgba(0, 0, 0, 0.5);

  /* Text colors */
  --color-text-primary: #333333;
  --color-text-secondary: #666666;
  --color-text-muted: #999999;
  --color-text-inverse: #ffffff;

  /* Accent colors */
  --color-accent: #667eea;
  --color-accent-light: #e7f3ff;
  --color-accent-text: #0056b3;

  /* Status colors - Success */
  --color-success: #4caf50;
  --color-success-bg: #d4edda;
  --color-success-text: #155724;

  /* Status colors - Warning */
  --color-warning: #ff9800;
  --color-warning-bg: #fff3cd;
  --color-warning-text: #856404;

  /* Status colors - Error */
  --color-error: #f44336;
  --color-error-bg: #f8d7da;
  --color-error-text: #721c24;

  /* Border colors */
  --color-border: #dee2e6;
  --color-border-light: rgba(255, 255, 255, 0.3);

  /* Interactive states */
  --color-hover-bg: rgba(255, 255, 255, 0.3);
  --color-focus-ring: rgba(255, 255, 255, 0.5);

  /* Shadows */
  --shadow-sm: 0 4px 6px rgba(0, 0, 0, 0.1);
  --shadow-md: 0 6px 16px rgba(0, 0, 0, 0.15);
  --shadow-lg: 0 20px 60px rgba(0, 0, 0, 0.3);

  /* Star banner specific */
  --color-star-gradient-start: rgba(255, 215, 0, 0.15);
  --color-star-gradient-end: rgba(255, 165, 0, 0.15);
  --color-star-border: rgba(255, 215, 0, 0.3);
  --color-star-highlight: #ffd700;
  --color-star-button-bg: linear-gradient(135deg, #ffd700, #ffed4e);
  --color-star-button-hover: linear-gradient(135deg, #ffed4e, #ffd700);
  --color-star-button-text: #333333;
  --color-star-button-shadow: rgba(255, 215, 0, 0.3);

  /* Code/monospace */
  --color-code-bg: #f8f9fa;
  --color-code-text: inherit;

  /* Transition */
  --theme-transition: background-color 0.2s ease, color 0.2s ease, border-color 0.2s ease, box-shadow 0.2s ease;
}

/* Dark Theme */
[data-theme="dark"] {
  /* Background colors */
  --color-bg-primary: #1e1e2e;
  --color-bg-secondary: #2a2a3e;
  --color-bg-tertiary: #3a3a4e;
  --color-bg-header: rgba(0, 0, 0, 0.3);
  --color-bg-overlay: rgba(0, 0, 0, 0.7);

  /* Text colors */
  --color-text-primary: #e4e4e7;
  --color-text-secondary: #a1a1aa;
  --color-text-muted: #71717a;
  --color-text-inverse: #e4e4e7;

  /* Accent colors */
  --color-accent: #818cf8;
  --color-accent-light: rgba(129, 140, 248, 0.2);
  --color-accent-text: #a5b4fc;

  /* Status colors - Success */
  --color-success: #4ade80;
  --color-success-bg: rgba(74, 222, 128, 0.2);
  --color-success-text: #86efac;

  /* Status colors - Warning */
  --color-warning: #fbbf24;
  --color-warning-bg: rgba(251, 191, 36, 0.2);
  --color-warning-text: #fcd34d;

  /* Status colors - Error */
  --color-error: #f87171;
  --color-error-bg: rgba(248, 113, 113, 0.2);
  --color-error-text: #fca5a5;

  /* Border colors */
  --color-border: #3f3f46;
  --color-border-light: rgba(255, 255, 255, 0.2);

  /* Interactive states */
  --color-hover-bg: rgba(255, 255, 255, 0.1);
  --color-focus-ring: rgba(255, 255, 255, 0.3);

  /* Shadows */
  --shadow-sm: 0 4px 6px rgba(0, 0, 0, 0.3);
  --shadow-md: 0 6px 16px rgba(0, 0, 0, 0.4);
  --shadow-lg: 0 20px 60px rgba(0, 0, 0, 0.5);

  /* Star banner specific */
  --color-star-gradient-start: rgba(255, 215, 0, 0.1);
  --color-star-gradient-end: rgba(255, 165, 0, 0.1);
  --color-star-border: rgba(255, 215, 0, 0.2);
  --color-star-highlight: #fcd34d;
  --color-star-button-bg: linear-gradient(135deg, #fbbf24, #fcd34d);
  --color-star-button-hover: linear-gradient(135deg, #fcd34d, #fbbf24);
  --color-star-button-text: #1e1e2e;
  --color-star-button-shadow: rgba(251, 191, 36, 0.3);

  /* Code/monospace */
  --color-code-bg: #2a2a3e;
  --color-code-text: #e4e4e7;
}

.app {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.header {
  background: var(--color-bg-header);
  backdrop-filter: blur(10px);
  padding: 2rem 0;
  color: var(--color-text-inverse);
  box-shadow: var(--shadow-sm);
  transition: var(--theme-transition);
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 2rem;
}

.title {
  font-size: 2.5rem;
  font-weight: 700;
  margin-bottom: 0.5rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.icon {
  font-size: 2rem;
}

.subtitle {
  font-size: 1.1rem;
  opacity: 0.9;
}

.header-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 2rem;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.refresh-button,
.theme-toggle {
  background: rgba(255, 255, 255, 0.2);
  border: 2px solid var(--color-border-light);
  color: var(--color-text-inverse);
  padding: 0.75rem 1.5rem;
  border-radius: 8px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  transition: all 0.3s;
  white-space: nowrap;
}

.theme-toggle {
  padding: 0.75rem;
  font-size: 1.25rem;
}

.refresh-button:hover:not(:disabled),
.theme-toggle:hover {
  background: var(--color-hover-bg);
  border-color: var(--color-focus-ring);
  transform: translateY(-2px);
}

.refresh-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.refresh-icon {
  font-size: 1.2rem;
  display: inline-block;
  transition: transform 0.3s;
}

.refresh-icon.spinning {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

.star-banner {
  background: linear-gradient(
    135deg,
    var(--color-star-gradient-start),
    var(--color-star-gradient-end)
  );
  backdrop-filter: blur(10px);
  border-bottom: 2px solid var(--color-star-border);
  animation: slideDown 0.5s ease-out;
  transition: var(--theme-transition);
}

@keyframes slideDown {
  from {
    opacity: 0;
    transform: translateY(-20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.star-banner-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1.5rem;
  padding: 1rem 0;
  flex-wrap: wrap;
}

.star-banner-text {
  display: flex;
  align-items: center;
  gap: 1rem;
  flex: 1;
  min-width: 300px;
}

.star-banner-icon {
  font-size: 2rem;
  animation: pulse 2s ease-in-out infinite;
  will-change: transform;
}

@keyframes pulse {
  0%,
  100% {
    transform: scale(1);
  }
  50% {
    transform: scale(1.1);
  }
}

.star-banner-message {
  color: var(--color-text-inverse);
  font-size: 0.95rem;
  line-height: 1.5;
}

.star-banner-message strong {
  font-weight: 600;
  color: var(--color-star-highlight);
}

.star-banner-actions {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.star-button {
  background: var(--color-star-button-bg);
  color: var(--color-star-button-text);
  padding: 0.65rem 1.5rem;
  border-radius: 8px;
  font-weight: 600;
  text-decoration: none;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  transition: all 0.3s;
  box-shadow: 0 4px 12px var(--color-star-button-shadow);
  white-space: nowrap;
}

.star-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px var(--color-star-button-shadow);
  background: var(--color-star-button-hover);
}

.star-button:focus-visible {
  outline: 2px solid var(--color-text-primary);
  outline-offset: 2px;
}
.dismiss-button {
  background: rgba(255, 255, 255, 0.2);
  border: 1px solid var(--color-border-light);
  color: var(--color-text-inverse);
  padding: 0.5rem;
  width: 2.5rem;
  height: 2.5rem;
  border-radius: 6px;
  font-size: 1.5rem;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s;
  line-height: 1;
}

.dismiss-button:hover {
  background: var(--color-hover-bg);
  border-color: var(--color-focus-ring);
}

.dismiss-button:focus-visible {
  outline: 2px solid var(--color-text-inverse);
  outline-offset: 2px;
}
.main {
  flex: 1;
  padding: 2rem 0;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.stat-card {
  background: var(--color-bg-primary);
  border-radius: 12px;
  padding: 1.5rem;
  display: flex;
  align-items: center;
  gap: 1rem;
  box-shadow: var(--shadow-sm);
  transition: transform 0.2s, var(--theme-transition);
}

.stat-card:hover {
  transform: translateY(-4px);
}

.stat-icon {
  font-size: 2.5rem;
}

.stat-content {
  flex: 1;
}

.stat-value {
  font-size: 2rem;
  font-weight: 700;
  color: var(--color-accent);
}

.stat-label {
  color: var(--color-text-secondary);
  font-size: 0.9rem;
}

.section {
  margin-bottom: 2rem;
}

.section-title {
  color: var(--color-text-inverse);
  font-size: 1.5rem;
  margin-bottom: 1rem;
  font-weight: 600;
}

.card {
  background: var(--color-bg-primary);
  border-radius: 12px;
  padding: 1.5rem;
  box-shadow: var(--shadow-sm);
  transition: var(--theme-transition);
}

.source-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.source-item {
  padding: 1rem;
  background: var(--color-bg-secondary);
  border-radius: 8px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 1rem;
  transition: var(--theme-transition);
}

.source-ip {
  font-family: "Courier New", monospace;
  font-weight: 600;
  color: var(--color-text-primary);
  min-width: 150px;
}

.source-stats {
  flex: 1;
}

.source-count {
  font-size: 0.9rem;
  color: var(--color-text-secondary);
  margin-bottom: 0.5rem;
}

.source-bar {
  height: 8px;
  background: var(--color-bg-tertiary);
  border-radius: 4px;
  overflow: hidden;
  display: flex;
  margin-bottom: 0.5rem;
}

.source-bar-pass {
  background: var(--color-success);
  transition: width 0.3s;
}

.source-bar-fail {
  background: var(--color-error);
  transition: width 0.3s;
}

.source-legend {
  display: flex;
  gap: 1rem;
  font-size: 0.85rem;
}

.legend-pass {
  color: var(--color-success);
}

.legend-fail {
  color: var(--color-error);
}

.table-container {
  overflow-x: auto;
}

.report-table {
  width: 100%;
  border-collapse: collapse;
}

.report-table th {
  text-align: left;
  padding: 0.75rem;
  background: var(--color-bg-secondary);
  font-weight: 600;
  color: var(--color-text-primary);
  border-bottom: 2px solid var(--color-border);
  transition: var(--theme-transition);
}

.report-table th.sortable {
  cursor: pointer;
  user-select: none;
  transition: background 0.2s;
}

.report-table th.sortable:hover {
  background: var(--color-bg-tertiary);
}

.report-table td {
  padding: 0.75rem;
  border-bottom: 1px solid var(--color-border);
  color: var(--color-text-primary);
  transition: var(--theme-transition);
}

.report-row {
  cursor: pointer;
  transition: background 0.2s;
}

.report-row:hover {
  background: var(--color-bg-secondary);
}

.date-cell {
  color: var(--color-text-secondary);
  font-size: 0.9rem;
}

code {
  background: var(--color-code-bg);
  color: var(--color-code-text);
  padding: 0.2rem 0.4rem;
  border-radius: 4px;
  font-family: "Courier New", monospace;
  font-size: 0.9rem;
  transition: var(--theme-transition);
}

.compliance-badge {
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  font-size: 0.85rem;
  font-weight: 600;
}

.compliance-badge.high {
  background: var(--color-success-bg);
  color: var(--color-success-text);
}

.compliance-badge.medium {
  background: var(--color-warning-bg);
  color: var(--color-warning-text);
}

.compliance-badge.low {
  background: var(--color-error-bg);
  color: var(--color-error-text);
}

.policy-badge {
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  background: var(--color-accent-light);
  color: var(--color-accent-text);
  font-size: 0.85rem;
  font-weight: 600;
  text-transform: uppercase;
}

.empty-state {
  text-align: center;
  padding: 3rem;
  color: var(--color-text-secondary);
}

.modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: var(--color-bg-overlay);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 2rem;
}

.modal-content {
  background: var(--color-bg-primary);
  border-radius: 12px;
  max-width: 800px;
  width: 100%;
  max-height: 80vh;
  overflow-y: auto;
  box-shadow: var(--shadow-lg);
  transition: var(--theme-transition);
}

.modal-header {
  padding: 1.5rem;
  border-bottom: 1px solid var(--color-border);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.modal-header h3 {
  margin: 0;
  font-size: 1.5rem;
  color: var(--color-text-primary);
}

.modal-close {
  background: none;
  border: none;
  font-size: 2rem;
  cursor: pointer;
  color: var(--color-text-muted);
  padding: 0;
  width: 2rem;
  height: 2rem;
  display: flex;
  align-items: center;
  justify-content: center;
}

.modal-close:hover {
  color: var(--color-text-primary);
}

.modal-body {
  padding: 1.5rem;
}

.detail-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
  margin-bottom: 1.5rem;
}

.detail-item {
  padding: 0.75rem;
  background: var(--color-bg-secondary);
  border-radius: 6px;
  color: var(--color-text-primary);
  transition: var(--theme-transition);
}

.detail-subtitle {
  margin: 1.5rem 0 1rem;
  font-size: 1.2rem;
  color: var(--color-text-primary);
}

.records-list {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.record-item {
  border: 1px solid var(--color-border);
  border-radius: 6px;
  padding: 1rem;
  transition: var(--theme-transition);
}

.record-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 0.75rem;
  font-weight: 600;
}

.record-ip {
  font-family: "Courier New", monospace;
  color: var(--color-accent);
}

.record-count {
  color: var(--color-text-secondary);
  font-size: 0.9rem;
}

.record-details {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.result-badge {
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  font-size: 0.85rem;
  background: var(--color-bg-tertiary);
  color: var(--color-text-secondary);
}

.result-badge.pass {
  background: var(--color-success-bg);
  color: var(--color-success-text);
}

.result-badge.fail {
  background: var(--color-error-bg);
  color: var(--color-error-text);
}

.footer {
  background: var(--color-bg-header);
  backdrop-filter: blur(10px);
  padding: 2rem 0 1rem;
  color: var(--color-text-inverse);
  margin-top: auto;
  transition: var(--theme-transition);
}

.footer-content {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 2rem;
  margin-bottom: 1.5rem;
}

.footer-section h4 {
  font-size: 1rem;
  font-weight: 600;
  margin-bottom: 0.75rem;
  color: rgba(255, 255, 255, 0.95);
}

.footer-section p {
  font-size: 0.85rem;
  color: rgba(255, 255, 255, 0.7);
  margin: 0;
}

.footer-links {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.footer-links a {
  color: rgba(255, 255, 255, 0.8);
  text-decoration: none;
  font-size: 0.85rem;
  transition: color 0.2s;
}

.footer-links a:hover {
  color: rgba(255, 255, 255, 1);
  text-decoration: underline;
}

.footer-bottom {
  border-top: 1px solid var(--color-border-light);
  padding-top: 1rem;
  text-align: center;
  font-size: 0.85rem;
  color: rgba(255, 255, 255, 0.7);
}

.footer-bottom p {
  margin: 0.25rem 0;
}

.footer-bottom a {
  color: rgba(255, 255, 255, 0.9);
  text-decoration: none;
}

.footer-bottom a:hover {
  text-decoration: underline;
}

.opensource-message {
  font-size: 0.8rem;
  color: rgba(255, 255, 255, 0.6);
  margin-top: 0.5rem;
}

.docker-image {
  background: rgba(0, 0, 0, 0.3);
  padding: 0.2rem 0.5rem;
  border-radius: 4px;
  font-family: "Courier New", monospace;
  font-size: 0.8rem;
  color: #88ddff;
}

@media (max-width: 768px) {
  .title {
    font-size: 1.75rem;
  }

  .header-top {
    flex-direction: column;
    align-items: flex-start;
    gap: 1rem;
  }

  .header-actions {
    width: 100%;
    justify-content: flex-end;
  }

  .refresh-button,
  .theme-toggle {
    flex: 1;
    justify-content: center;
  }

  .star-banner-content {
    flex-direction: column;
    align-items: stretch;
    gap: 1rem;
  }

  .star-banner-text {
    min-width: auto;
  }

  .star-banner-actions {
    width: 100%;
    justify-content: space-between;
  }

  .star-button {
    flex: 1;
    justify-content: center;
  }

  .stats-grid {
    grid-template-columns: 1fr;
  }

  .source-item {
    flex-direction: column;
    align-items: flex-start;
  }

  .source-ip {
    min-width: auto;
  }

  .footer-content {
    grid-template-columns: 1fr;
    gap: 1.5rem;
  }

  .footer-section {
    text-align: center;
  }

  .footer-links {
    align-items: center;
  }
}
</style>
