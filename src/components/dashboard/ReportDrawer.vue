<script setup>
import { computed, ref } from "vue";

const props = defineProps({
  isOpen: Boolean,
  report: {
    type: Object,
    default: null,
  },
  loading: {
    type: Boolean,
    default: false,
  },
});

const emit = defineEmits(["close"]);

const activeTab = ref("overview");

// Close on backdrop click
const handleBackdropClick = () => {
  emit("close");
};

// Close on escape key
const handleKeydown = (e) => {
  if (e.key === "Escape") emit("close");
};

// Helpers for status styling
const getResultColor = (result) => {
  const r = result?.toLowerCase();
  if (r === "pass") return "text-success";
  if (r === "fail" || r === "hardfail") return "text-danger";
  return "text-warning";
};

const getBadgeClass = (result) => {
  const r = result?.toLowerCase();
  if (r === "pass") return "badge-success";
  if (r === "fail" || r === "hardfail") return "badge-danger";
  return "badge-warning";
};

// Compute overall DMARC result from records
const dmarcResult = computed(() => {
  if (!props.report?.Records?.length) return "unknown";
  const hasPass = props.report.Records.some(
    (r) =>
      r.Row?.PolicyEvaluated?.DKIM === "pass" ||
      r.Row?.PolicyEvaluated?.SPF === "pass",
  );
  return hasPass ? "pass" : "fail";
});

// Get total message count
const totalMessages = computed(() => {
  if (!props.report?.Records) return 0;
  return props.report.Records.reduce((sum, r) => sum + (r.Row?.Count || 0), 0);
});

// Format date
const formatDate = (timestamp) => {
  if (!timestamp) return "N/A";
  const date = new Date(timestamp * 1000);
  return new Intl.DateTimeFormat("en-US", {
    year: "numeric",
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  }).format(date);
};

// Check alignment - simple check if domain matches header_from
const isAligned = (authDomain, record) => {
  if (!authDomain || !record?.Identifiers?.HeaderFrom) return false;
  const headerFrom = record.Identifiers.HeaderFrom.toLowerCase();
  return (
    authDomain.toLowerCase().endsWith(headerFrom) ||
    headerFrom.endsWith(authDomain.toLowerCase())
  );
};
</script>

<template>
  <Teleport to="body">
    <div
      class="drawer-root"
      :class="{ 'is-open': isOpen }"
      @keydown="handleKeydown"
    >
      <div class="backdrop" @click="handleBackdropClick"></div>

      <div class="panel" @click.stop>
        <!-- Header -->
        <header class="panel-header" v-if="report">
          <div>
            <h2 class="panel-title">
              {{ report.ReportMetadata?.OrgName || "Report Details" }}
            </h2>
            <div class="panel-subtitle">
              {{ report.PolicyPublished?.Domain || "Unknown Domain" }}
            </div>
          </div>
          <button
            class="btn-close"
            @click="emit('close')"
            aria-label="Close report details panel"
          >
            <svg
              width="24"
              height="24"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
            >
              <line x1="18" y1="6" x2="6" y2="18"></line>
              <line x1="6" y1="6" x2="18" y2="18"></line>
            </svg>
          </button>
        </header>

        <!-- Loading State -->
        <div v-if="loading" class="panel-loading">
          <div class="spinner"></div>
          <p>Loading report details...</p>
        </div>

        <!-- Body -->
        <div class="panel-body" v-else-if="report">
          <!-- DMARC Verdict Card -->
          <div class="verdict-card" :class="getBadgeClass(dmarcResult)">
            <div class="verdict-row">
              <span class="verdict-label">DMARC Result</span>
              <span class="verdict-value">{{ dmarcResult.toUpperCase() }}</span>
            </div>
            <div class="verdict-meta">
              <span
                >Policy:
                <strong>{{ report.PolicyPublished?.P || "none" }}</strong></span
              >
              <span
                >Volume:
                <strong
                  >{{ totalMessages.toLocaleString() }} emails</strong
                ></span
              >
            </div>
          </div>

          <!-- Tabs -->
          <div class="tabs">
            <button
              :class="{ active: activeTab === 'overview' }"
              @click="activeTab = 'overview'"
            >
              Overview
            </button>
            <button
              :class="{ active: activeTab === 'records' }"
              @click="activeTab = 'records'"
            >
              Records ({{ report.Records?.length || 0 }})
            </button>
            <button
              :class="{ active: activeTab === 'raw' }"
              @click="activeTab = 'raw'"
            >
              Raw Data
            </button>
          </div>

          <!-- Overview Tab -->
          <div v-if="activeTab === 'overview'" class="tab-content">
            <div class="info-grid">
              <div class="info-item">
                <span class="info-label">Report ID</span>
                <span class="info-value font-mono">{{
                  report.ReportMetadata?.ReportID || "N/A"
                }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">Contact Email</span>
                <span class="info-value">{{
                  report.ReportMetadata?.Email || "N/A"
                }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">Date Range</span>
                <span class="info-value">
                  {{ formatDate(report.ReportMetadata?.DateRange?.Begin) }} -
                  {{ formatDate(report.ReportMetadata?.DateRange?.End) }}
                </span>
              </div>
              <div class="info-item">
                <span class="info-label">DKIM Alignment</span>
                <span class="info-value">{{
                  report.PolicyPublished?.ADKIM === "s" ? "Strict" : "Relaxed"
                }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">SPF Alignment</span>
                <span class="info-value">{{
                  report.PolicyPublished?.ASPF === "s" ? "Strict" : "Relaxed"
                }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">Subdomain Policy</span>
                <span class="info-value">{{
                  report.PolicyPublished?.SP || "Same as domain"
                }}</span>
              </div>
            </div>
          </div>

          <!-- Records Tab -->
          <div v-if="activeTab === 'records'" class="tab-content">
            <div v-if="!report.Records?.length" class="empty-state">
              No records found in this report.
            </div>

            <div
              v-for="(record, idx) in report.Records || []"
              :key="idx"
              class="record-card"
            >
              <div class="record-header">
                <div class="record-ip">
                  <svg
                    width="16"
                    height="16"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                  >
                    <rect x="2" y="2" width="20" height="8" rx="2" ry="2" />
                    <rect x="2" y="14" width="20" height="8" rx="2" ry="2" />
                    <line x1="6" y1="6" x2="6.01" y2="6" />
                    <line x1="6" y1="18" x2="6.01" y2="18" />
                  </svg>
                  <span class="font-mono">{{
                    record.Row?.SourceIP || "N/A"
                  }}</span>
                </div>
                <span class="record-count"
                  >{{ record.Row?.Count || 0 }} messages</span
                >
              </div>

              <div class="auth-grid">
                <!-- SPF Section -->
                <section class="auth-col">
                  <h4 class="section-heading">
                    <svg
                      width="14"
                      height="14"
                      viewBox="0 0 24 24"
                      fill="none"
                      stroke="currentColor"
                      stroke-width="2"
                    >
                      <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z" />
                    </svg>
                    SPF
                    <span
                      class="status-pill"
                      :class="getResultColor(record.Row?.PolicyEvaluated?.SPF)"
                    >
                      {{ record.Row?.PolicyEvaluated?.SPF || "unknown" }}
                    </span>
                  </h4>

                  <div
                    v-for="(spf, spfIdx) in record.AuthResults?.SPF || []"
                    :key="'spf-' + spfIdx"
                    class="auth-item"
                  >
                    <div class="auth-header">
                      <span class="domain-mono">{{ spf.Domain }}</span>
                      <span
                        class="status-text"
                        :class="getResultColor(spf.Result)"
                      >
                        {{ spf.Result?.toUpperCase() }}
                      </span>
                    </div>
                    <div class="auth-details">
                      <span>Scope: {{ spf.Scope || "mfrom" }}</span>
                      <span
                        v-if="isAligned(spf.Domain, record)"
                        class="alignment-tag aligned"
                        >Aligned</span
                      >
                      <span v-else class="alignment-tag unaligned"
                        >Unaligned</span
                      >
                    </div>
                  </div>

                  <div
                    v-if="!record.AuthResults?.SPF?.length"
                    class="empty-auth"
                  >
                    No SPF results
                  </div>
                </section>

                <!-- DKIM Section -->
                <section class="auth-col">
                  <h4 class="section-heading">
                    <svg
                      width="14"
                      height="14"
                      viewBox="0 0 24 24"
                      fill="none"
                      stroke="currentColor"
                      stroke-width="2"
                    >
                      <rect x="3" y="11" width="18" height="11" rx="2" ry="2" />
                      <path d="M7 11V7a5 5 0 0 1 10 0v4" />
                    </svg>
                    DKIM
                    <span
                      class="status-pill"
                      :class="getResultColor(record.Row?.PolicyEvaluated?.DKIM)"
                    >
                      {{ record.Row?.PolicyEvaluated?.DKIM || "unknown" }}
                    </span>
                  </h4>

                  <div
                    v-for="(dkim, dkimIdx) in record.AuthResults?.DKIM || []"
                    :key="'dkim-' + dkimIdx"
                    class="auth-item"
                  >
                    <div class="auth-header">
                      <span class="domain-mono">{{ dkim.Domain }}</span>
                      <span
                        class="status-text"
                        :class="getResultColor(dkim.Result)"
                      >
                        {{ dkim.Result?.toUpperCase() }}
                      </span>
                    </div>
                    <div class="auth-details">
                      <span>Selector: {{ dkim.Selector || "N/A" }}</span>
                      <span
                        v-if="isAligned(dkim.Domain, record)"
                        class="alignment-tag aligned"
                        >Aligned</span
                      >
                      <span v-else class="alignment-tag unaligned"
                        >Unaligned</span
                      >
                    </div>
                  </div>

                  <div
                    v-if="!record.AuthResults?.DKIM?.length"
                    class="empty-auth"
                  >
                    No DKIM signatures
                  </div>
                </section>
              </div>

              <div class="record-footer">
                <span
                  >Disposition:
                  <strong>{{
                    record.Row?.PolicyEvaluated?.Disposition || "none"
                  }}</strong></span
                >
                <span v-if="record.Identifiers?.HeaderFrom">
                  From:
                  <span class="font-mono">{{
                    record.Identifiers.HeaderFrom
                  }}</span>
                </span>
              </div>
            </div>
          </div>

          <!-- Raw Data Tab -->
          <div v-if="activeTab === 'raw'" class="tab-content">
            <div class="raw-data-section">
              <pre class="code-block">{{
                JSON.stringify(report, null, 2)
              }}</pre>
            </div>
          </div>
        </div>

        <!-- Empty State -->
        <div v-else class="panel-empty">
          <p>Select a report to view details</p>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<style scoped>
/* --- Drawer Mechanism --- */
.drawer-root {
  position: fixed;
  inset: 0;
  z-index: 50;
  pointer-events: none;
  visibility: hidden;
  transition: visibility 0.3s;
}

.drawer-root.is-open {
  pointer-events: auto;
  visibility: visible;
}

.backdrop {
  position: absolute;
  inset: 0;
  background: rgba(15, 23, 42, 0.5);
  backdrop-filter: blur(2px);
  opacity: 0;
  transition: opacity 0.3s;
}

.is-open .backdrop {
  opacity: 1;
}

.panel {
  position: absolute;
  top: 0;
  right: 0;
  bottom: 0;
  width: 100%;
  max-width: 540px;
  background: var(--bg-card);
  border-left: 1px solid var(--border-subtle);
  box-shadow: -4px 0 24px rgba(0, 0, 0, 0.15);
  transform: translateX(100%);
  transition: transform 0.3s cubic-bezier(0.16, 1, 0.3, 1);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.is-open .panel {
  transform: translateX(0);
}

/* --- Panel Header --- */
.panel-header {
  padding: 20px 24px;
  border-bottom: 1px solid var(--border-subtle);
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  background: var(--bg-app);
  flex-shrink: 0;
}

.panel-title {
  margin: 0;
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--text-main);
}

.panel-subtitle {
  color: var(--text-muted);
  font-size: 0.875rem;
  font-family: var(--font-mono);
  margin-top: 4px;
}

.btn-close {
  background: none;
  border: none;
  cursor: pointer;
  color: var(--text-muted);
  padding: 4px;
  border-radius: 4px;
  transition: background 0.2s;
  flex-shrink: 0;
}
.btn-close:hover {
  background: var(--border-subtle);
  color: var(--text-main);
}

/* --- Loading State --- */
.panel-loading {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 16px;
  color: var(--text-muted);
}

.spinner {
  width: 32px;
  height: 32px;
  border: 3px solid var(--border-subtle);
  border-top-color: var(--c-primary);
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

/* --- Panel Body --- */
.panel-body {
  padding: 20px 24px;
  overflow-y: auto;
  flex: 1;
}

.panel-empty {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-muted);
}

/* --- Verdict Card --- */
.verdict-card {
  padding: 16px;
  border-radius: var(--radius-card);
  margin-bottom: 20px;
  border: 1px solid transparent;
}

.badge-success {
  background: var(--c-success-soft);
  border-color: var(--c-success);
  color: var(--c-success);
}
.badge-warning {
  background: var(--c-warning-soft);
  border-color: var(--c-warning);
  color: var(--c-warning);
}
.badge-danger {
  background: var(--c-danger-soft);
  border-color: var(--c-danger);
  color: var(--c-danger);
}

.verdict-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}
.verdict-label {
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  opacity: 0.8;
}
.verdict-value {
  font-size: 1.25rem;
  font-weight: 800;
}

.verdict-meta {
  font-size: 0.875rem;
  opacity: 0.9;
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
}

/* --- Tabs --- */
.tabs {
  display: flex;
  gap: 4px;
  border-bottom: 1px solid var(--border-subtle);
  margin-bottom: 20px;
}

.tabs button {
  background: transparent;
  border: none;
  padding: 10px 16px;
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--text-muted);
  cursor: pointer;
  border-bottom: 2px solid transparent;
  margin-bottom: -1px;
  transition: all 0.2s;
  font-family: var(--font-sans);
}

.tabs button:hover {
  color: var(--text-main);
}

.tabs button.active {
  color: var(--c-primary);
  border-bottom-color: var(--c-primary);
}

/* --- Tab Content --- */
.tab-content {
  animation: fadeIn 0.2s ease;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(4px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* --- Info Grid (Overview) --- */
.info-grid {
  display: grid;
  gap: 16px;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.info-label {
  font-size: 0.75rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: var(--text-muted);
  font-weight: 600;
}

.info-value {
  font-size: 0.875rem;
  color: var(--text-main);
}

/* --- Record Card --- */
.record-card {
  background: var(--bg-app);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-card);
  margin-bottom: 16px;
  overflow: hidden;
}

.record-header {
  padding: 12px 16px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid var(--border-subtle);
  background: var(--bg-card);
}

.record-ip {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--text-main);
  font-weight: 500;
}

.record-ip svg {
  color: var(--text-muted);
}

.record-count {
  font-size: 0.8125rem;
  color: var(--text-muted);
  font-family: var(--font-mono);
}

/* --- Auth Grid --- */
.auth-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1px;
  background: var(--border-subtle);
}

.auth-col {
  padding: 12px 16px;
  background: var(--bg-app);
}

.section-heading {
  font-size: 0.75rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: var(--text-muted);
  margin: 0 0 12px 0;
  display: flex;
  align-items: center;
  gap: 6px;
  font-weight: 600;
}

.status-pill {
  font-size: 0.65rem;
  padding: 2px 6px;
  border-radius: 4px;
  text-transform: uppercase;
  margin-left: auto;
  background: var(--bg-card);
}

.auth-item {
  background: var(--bg-card);
  border: 1px solid var(--border-subtle);
  border-radius: 6px;
  padding: 10px;
  margin-bottom: 8px;
}

.auth-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 4px;
}

.domain-mono {
  font-family: var(--font-mono);
  font-size: 0.75rem;
  color: var(--text-main);
  word-break: break-all;
}

.status-text {
  font-weight: 700;
  font-size: 0.7rem;
  text-transform: uppercase;
}
.text-success {
  color: var(--c-success);
}
.text-warning {
  color: var(--c-warning);
}
.text-danger {
  color: var(--c-danger);
}

.auth-details {
  font-size: 0.7rem;
  color: var(--text-muted);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.alignment-tag {
  font-size: 0.6rem;
  padding: 2px 6px;
  border-radius: 4px;
  text-transform: uppercase;
  font-weight: 600;
}
.aligned {
  background: var(--c-success-soft);
  color: var(--c-success);
}
.unaligned {
  background: var(--border-subtle);
  color: var(--text-muted);
}

.empty-auth {
  font-size: 0.8125rem;
  color: var(--text-muted);
  font-style: italic;
}

.record-footer {
  padding: 10px 16px;
  font-size: 0.8125rem;
  color: var(--text-muted);
  display: flex;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 8px;
  border-top: 1px solid var(--border-subtle);
}

/* --- Raw Data --- */
.raw-data-section {
  overflow: hidden;
  border-radius: var(--radius-card);
}

.code-block {
  background: var(--code-bg);
  color: #d4d4d4;
  padding: 16px;
  font-family: var(--font-mono);
  font-size: 0.75rem;
  overflow-x: auto;
  white-space: pre-wrap;
  word-break: break-all;
  margin: 0;
  max-height: 60vh;
  border-radius: var(--radius-card);
}

/* --- Empty State --- */
.empty-state {
  text-align: center;
  color: var(--text-muted);
  padding: 32px;
  font-style: italic;
}

/* --- Utility --- */
.font-mono {
  font-family: var(--font-mono);
}

/* --- Responsive --- */
@media (max-width: 600px) {
  .panel {
    max-width: 100%;
  }

  .auth-grid {
    grid-template-columns: 1fr;
  }

  .verdict-meta {
    flex-direction: column;
    gap: 4px;
  }
}
</style>
