<script setup>
import { computed } from "vue";

const props = defineProps({
  complianceScore: { type: Number, required: true },
  volume: { type: Number, default: 0 },
  sourceCount: { type: Number, default: 0 },
  dateRange: { type: String, default: "Last 30 Days" },
  loading: { type: Boolean, default: false },
});

const emit = defineEmits(["refresh"]);

// Determine health state based on score
const healthState = computed(() => {
  if (props.complianceScore >= 95) return "secure";
  if (props.complianceScore >= 80) return "warning";
  return "critical";
});

// Format numbers (e.g. 12,450)
const fmt = (n) => new Intl.NumberFormat("en-US").format(n);

// Dynamic status text
const statusMessage = computed(() => {
  if (healthState.value === "secure") return "Domain Protected";
  if (healthState.value === "warning") return "Compliance Issues Detected";
  return "Critical Authentication Gaps";
});

const statusSubtext = computed(() => {
  if (healthState.value === "secure") return "Traffic is fully authenticated.";
  if (healthState.value === "warning")
    return "Some legitimate email may be failing.";
  return "High risk of spoofing. Immediate action required.";
});
</script>

<template>
  <div class="hero-container">
    <header class="page-header">
      <div>
        <h1 class="page-title">Dashboard</h1>
        <p class="page-subtitle">Overview for {{ dateRange }}</p>
      </div>
      <div class="actions">
        <button class="btn-ghost" @click="emit('refresh')" :disabled="loading">
          <svg
            :class="{ spinning: loading }"
            width="16"
            height="16"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
          >
            <path d="M21 12a9 9 0 0 0-9-9 9.75 9.75 0 0 0-6.74 2.74L3 8" />
            <path d="M3 3v5h5" />
            <path d="M3 12a9 9 0 0 0 9 9 9.75 9.75 0 0 0 6.74-2.74L21 16" />
            <path d="M16 21h5v-5" />
          </svg>
          Refresh
        </button>
      </div>
    </header>

    <div class="stat-grid">
      <div class="card health-card" :class="healthState">
        <div class="health-content">
          <div class="icon-wrapper">
            <svg
              v-if="healthState === 'secure'"
              width="24"
              height="24"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z" />
              <path d="m9 12 2 2 4-4" />
            </svg>
            <svg
              v-else
              width="24"
              height="24"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z" />
              <path d="M12 8v4" />
              <path d="M12 16h.01" />
            </svg>
          </div>
          <div>
            <div class="stat-label">System Health</div>
            <div class="stat-value">{{ statusMessage }}</div>
            <div class="stat-desc">{{ statusSubtext }}</div>
          </div>
        </div>

        <div class="health-score">
          <span class="score-val">{{ complianceScore.toFixed(1) }}%</span>
          <span class="score-label">Pass Rate</span>
        </div>
      </div>

      <div class="card metric-card">
        <div class="metric-header">
          <span class="stat-label">Total Volume</span>
          <svg
            class="icon-subtle"
            width="20"
            height="20"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
          >
            <polyline points="22 12 18 12 15 21 9 3 6 12 2 12" />
          </svg>
        </div>
        <div class="stat-value mono">{{ fmt(volume) }}</div>
        <div class="stat-desc">Emails processed</div>
      </div>

      <div class="card metric-card">
        <div class="metric-header">
          <span class="stat-label">Sending Sources</span>
          <svg
            class="icon-subtle"
            width="20"
            height="20"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
          >
            <rect x="2" y="2" width="20" height="8" rx="2" ry="2" />
            <rect x="2" y="14" width="20" height="8" rx="2" ry="2" />
            <line x1="6" y1="6" x2="6.01" y2="6" />
            <line x1="6" y1="18" x2="6.01" y2="18" />
          </svg>
        </div>
        <div class="stat-value mono">{{ fmt(sourceCount) }}</div>
        <div class="stat-desc">Unique IPs detected</div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* --- Layout --- */
.hero-container {
  display: flex;
  flex-direction: column;
  gap: 24px;
  width: 100%;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  padding-bottom: 8px;
  border-bottom: 1px solid var(--border-subtle);
}

.page-title {
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--text-main);
  margin: 0;
  letter-spacing: -0.02em;
}

.page-subtitle {
  color: var(--text-muted);
  font-size: 0.875rem;
  margin: 4px 0 0 0;
}

.actions {
  display: flex;
  gap: 8px;
}

/* --- Grid System --- */
.stat-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 16px;
}

/* Tablet & Up */
@media (min-width: 768px) {
  .stat-grid {
    grid-template-columns: 2fr 1fr 1fr;
  }
}

/* --- Card Base --- */
.card {
  background: var(--bg-card);
  border: 1px solid var(--border-subtle);
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
}

/* --- Text Utility --- */
.stat-label {
  font-size: 0.75rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  font-weight: 600;
  color: var(--text-muted);
  margin-bottom: 4px;
}

.stat-value {
  font-size: 1.5rem;
  font-weight: 600;
  color: var(--text-main);
  line-height: 1.2;
}

.stat-value.mono {
  font-family: var(--font-mono);
  letter-spacing: -0.03em;
}

.stat-desc {
  font-size: 0.875rem;
  color: var(--text-muted);
  margin-top: 4px;
}

/* --- Health Card Styling --- */
.health-card {
  display: flex;
  justify-content: space-between;
  align-items: center;
  position: relative;
  overflow: hidden;
}

.health-content {
  display: flex;
  gap: 16px;
  align-items: flex-start;
}

.icon-wrapper {
  padding: 10px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.health-score {
  text-align: right;
  display: flex;
  flex-direction: column;
  align-items: flex-end;
}

.score-val {
  font-size: 2rem;
  font-weight: 700;
  font-family: var(--font-mono);
}

.score-label {
  font-size: 0.75rem;
  color: var(--text-muted);
  text-transform: uppercase;
  font-weight: 600;
}

/* --- Status Variants --- */

/* Secure (Green) */
.health-card.secure {
  border-left: 4px solid var(--c-success);
}
.health-card.secure .icon-wrapper {
  background-color: var(--c-success-soft);
  color: var(--c-success);
}
.health-card.secure .score-val {
  color: var(--c-success);
}

/* Warning (Amber) */
.health-card.warning {
  border-left: 4px solid var(--c-warning);
}
.health-card.warning .icon-wrapper {
  background-color: var(--c-warning-soft);
  color: var(--c-warning);
}
.health-card.warning .score-val {
  color: var(--c-warning);
}

/* Critical (Red) */
.health-card.critical {
  border-left: 4px solid var(--c-danger);
}
.health-card.critical .icon-wrapper {
  background-color: var(--c-danger-soft);
  color: var(--c-danger);
}
.health-card.critical .score-val {
  color: var(--c-danger);
}

/* --- Metric Cards --- */
.metric-card {
  display: flex;
  flex-direction: column;
  justify-content: center;
}
.metric-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}
.icon-subtle {
  color: var(--text-muted);
  opacity: 0.5;
}

/* --- Buttons --- */
.btn-ghost {
  background: transparent;
  border: none;
  color: var(--text-muted);
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 0.875rem;
  padding: 6px 12px;
  border-radius: 6px;
  transition: all 0.2s;
  font-family: var(--font-sans);
}
.btn-ghost:hover {
  background: var(--border-subtle);
  color: var(--text-main);
}
.btn-ghost:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* Spinning animation for refresh */
@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}
.spinning {
  animation: spin 1s linear infinite;
}

/* --- Responsive --- */
@media (max-width: 767px) {
  .health-card {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }

  .health-score {
    align-items: flex-start;
    width: 100%;
    padding-top: 16px;
    border-top: 1px solid var(--border-subtle);
  }

  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }

  .actions {
    width: 100%;
  }

  .btn-ghost {
    flex: 1;
    justify-content: center;
  }
}
</style>
