<script setup>
import { ref, computed } from "vue";

const props = defineProps({
  reports: {
    type: Array,
    default: () => [],
  },
});

const emit = defineEmits(["view-details"]);

const searchQuery = ref("");

// Filter reports based on search query
const filteredReports = computed(() => {
  if (!searchQuery.value) return props.reports;
  const query = searchQuery.value.toLowerCase();
  return props.reports.filter(
    (report) =>
      report.org_name?.toLowerCase().includes(query) ||
      report.domain?.toLowerCase().includes(query),
  );
});

// Helper to format dates
const formatDate = (timestamp) => {
  if (!timestamp) return "N/A";
  // Handle unix timestamp (seconds)
  const date = new Date(timestamp * 1000);
  return new Intl.DateTimeFormat("en-US", {
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  }).format(date);
};

// Helper for status badge styling
const getStatusClass = (report) => {
  const rate = report.compliance_rate ?? 0;
  if (rate >= 100) return "status-pass";
  if (rate >= 80) return "status-warn";
  return "status-fail";
};

const getStatusLabel = (report) => {
  const rate = report.compliance_rate ?? 0;
  if (rate >= 100) return "PASS";
  if (rate >= 80) return "WARN";
  return "FAIL";
};

// Get policy badge class
const getPolicyClass = (policy) => {
  if (policy === "reject") return "policy-reject";
  if (policy === "quarantine") return "policy-quarantine";
  return "policy-none";
};
</script>

<template>
  <div class="card table-container">
    <div class="table-toolbar">
      <h2 class="section-title">Recent Reports</h2>

      <div class="search-wrapper">
        <svg
          class="search-icon"
          width="16"
          height="16"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
        >
          <circle cx="11" cy="11" r="8" />
          <path d="m21 21-4.3-4.3" />
        </svg>
        <input
          type="text"
          v-model="searchQuery"
          placeholder="Filter domain or org..."
          class="search-input"
          aria-label="Filter domain or organization"
        />
      </div>
    </div>

    <div class="table-responsive">
      <table class="report-table">
        <thead>
          <tr>
            <th class="col-org">Organization</th>
            <th class="col-domain">Domain</th>
            <th class="col-vol text-right">Messages</th>
            <th class="col-rate text-right">Compliance</th>
            <th class="col-policy">Policy</th>
            <th class="col-date">Date Range</th>
            <th class="col-status">Status</th>
            <th class="col-action"></th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="report in filteredReports"
            :key="report.id"
            @click="emit('view-details', report)"
            class="clickable-row"
          >
            <td class="col-org">
              <div class="org-cell">
                <div class="org-icon">
                  <svg
                    width="16"
                    height="16"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                  >
                    <path d="M3 21h18" />
                    <path d="M5 21V7l8-4 8 4v14" />
                    <path d="M17 21v-5a2 2 0 0 0-2-2h-6a2 2 0 0 0-2 2v5" />
                  </svg>
                </div>
                <span class="org-name">{{ report.org_name || "Unknown" }}</span>
              </div>
            </td>

            <td class="col-domain font-mono">
              {{ report.domain }}
            </td>

            <td class="col-vol text-right font-mono">
              {{ report.total_messages?.toLocaleString() || 0 }}
            </td>

            <td class="col-rate text-right font-mono">
              {{ (report.compliance_rate ?? 0).toFixed(1) }}%
            </td>

            <td class="col-policy">
              <span
                class="policy-badge"
                :class="getPolicyClass(report.policy_p)"
              >
                {{ report.policy_p || "none" }}
              </span>
            </td>

            <td class="col-date text-muted">
              {{ formatDate(report.date_begin) }}
            </td>

            <td class="col-status">
              <span class="badge" :class="getStatusClass(report)">
                <span class="dot"></span>
                {{ getStatusLabel(report) }}
              </span>
            </td>

            <td class="col-action">
              <svg
                class="arrow-icon"
                width="16"
                height="16"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
              >
                <polyline points="9 18 15 12 9 6" />
              </svg>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-if="filteredReports.length === 0" class="empty-table">
      <svg
        width="48"
        height="48"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="1.5"
      >
        <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" />
        <polyline points="14 2 14 8 20 8" />
        <line x1="12" y1="18" x2="12" y2="12" />
        <line x1="9" y1="15" x2="15" y2="15" />
      </svg>
      <p v-if="searchQuery">No reports matching "{{ searchQuery }}"</p>
      <p v-else>No reports found.</p>
    </div>
  </div>
</template>

<style scoped>
/* --- Container --- */
.table-container {
  background: var(--bg-card);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-card);
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.table-toolbar {
  padding: 16px 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid var(--border-subtle);
  flex-wrap: wrap;
  gap: 12px;
}

.section-title {
  font-size: 1rem;
  font-weight: 600;
  margin: 0;
  color: var(--text-main);
}

/* --- Search Input --- */
.search-wrapper {
  position: relative;
  width: 240px;
  max-width: 100%;
}
.search-icon {
  position: absolute;
  left: 10px;
  top: 50%;
  transform: translateY(-50%);
  color: var(--text-muted);
  pointer-events: none;
}
.search-input {
  width: 100%;
  padding: 8px 12px 8px 32px;
  border: 1px solid var(--border-subtle);
  border-radius: 6px;
  font-family: var(--font-sans);
  font-size: 0.875rem;
  background: var(--bg-app);
  color: var(--text-main);
  transition: border-color 0.2s;
}
.search-input::placeholder {
  color: var(--text-muted);
}
.search-input:focus {
  outline: none;
  border-color: var(--c-primary);
}

/* --- Table Structure --- */
.table-responsive {
  width: 100%;
  overflow-x: auto;
}

.report-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.875rem;
  white-space: nowrap;
}

th {
  text-align: left;
  padding: 12px 20px;
  background: var(--bg-app);
  color: var(--text-muted);
  font-weight: 600;
  font-size: 0.75rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  border-bottom: 1px solid var(--border-subtle);
}

td {
  padding: 14px 20px;
  border-bottom: 1px solid var(--border-subtle);
  color: var(--text-main);
  transition: background 0.15s;
}

/* --- Row Interaction --- */
.clickable-row {
  cursor: pointer;
}
.clickable-row:hover {
  background-color: var(--bg-app);
}
.clickable-row:hover .arrow-icon {
  transform: translateX(4px);
  color: var(--c-primary);
}
.clickable-row:last-child td {
  border-bottom: none;
}

/* --- Column Specifics --- */
.org-cell {
  display: flex;
  align-items: center;
  gap: 10px;
}
.org-icon {
  width: 28px;
  height: 28px;
  background: var(--bg-app);
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-muted);
  flex-shrink: 0;
}
.org-name {
  font-weight: 500;
}

.font-mono {
  font-family: var(--font-mono);
  font-size: 0.8125rem;
}

.text-muted {
  color: var(--text-muted);
}

.text-right {
  text-align: right;
}

.arrow-icon {
  color: var(--border-subtle);
  transition:
    transform 0.2s,
    color 0.2s;
}

/* --- Policy Badge --- */
.policy-badge {
  display: inline-block;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 0.7rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.02em;
}

.policy-none {
  background: var(--bg-app);
  color: var(--text-muted);
  border: 1px solid var(--border-subtle);
}

.policy-quarantine {
  background: var(--c-warning-soft);
  color: var(--c-warning);
}

.policy-reject {
  background: var(--c-danger-soft);
  color: var(--c-danger);
}

/* --- Status Badges --- */
.badge {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 4px 10px;
  border-radius: 99px;
  font-size: 0.75rem;
  font-weight: 600;
  letter-spacing: 0.02em;
}

.dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background-color: currentColor;
}

/* Pass (Green) */
.status-pass {
  background-color: var(--c-success-soft);
  color: var(--c-success);
}

/* Warn (Amber) */
.status-warn {
  background-color: var(--c-warning-soft);
  color: var(--c-warning);
}

/* Fail (Red) */
.status-fail {
  background-color: var(--c-danger-soft);
  color: var(--c-danger);
}

/* --- Empty State --- */
.empty-table {
  padding: 48px 20px;
  text-align: center;
  color: var(--text-muted);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}

.empty-table svg {
  opacity: 0.3;
}

.empty-table p {
  margin: 0;
  font-size: 0.875rem;
}

/* --- Responsive --- */
@media (max-width: 768px) {
  .table-toolbar {
    flex-direction: column;
    align-items: stretch;
  }

  .search-wrapper {
    width: 100%;
  }

  th,
  td {
    padding: 12px 16px;
  }
}
</style>
