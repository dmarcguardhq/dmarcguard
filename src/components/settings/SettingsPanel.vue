<script setup>
import { ref, onMounted, watch } from "vue";

const props = defineProps({
  isOpen: {
    type: Boolean,
    default: false,
  },
});

const emit = defineEmits(["close"]);

const settings = ref({
  mcp_enabled: false,
  mcp_path: "/mcp",
});
const loading = ref(false);
const saving = ref(false);
const error = ref(null);
const success = ref(null);

const fetchSettings = async () => {
  loading.value = true;
  error.value = null;
  try {
    const response = await fetch("/api/settings");
    if (!response.ok) {
      throw new Error("Failed to fetch settings");
    }
    settings.value = await response.json();
  } catch (err) {
    error.value = err.message;
  } finally {
    loading.value = false;
  }
};

const toggleMCP = async () => {
  saving.value = true;
  error.value = null;
  success.value = null;

  const newValue = !settings.value.mcp_enabled;

  try {
    const response = await fetch("/api/settings/mcp", {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ enabled: newValue }),
    });

    if (!response.ok) {
      const data = await response.json();
      throw new Error(data.error || "Failed to update setting");
    }

    const data = await response.json();
    settings.value.mcp_enabled = data.enabled;
    success.value = data.message;

    // Clear success message after 3 seconds
    setTimeout(() => {
      success.value = null;
    }, 3000);
  } catch (err) {
    error.value = err.message;
  } finally {
    saving.value = false;
  }
};

const closePanel = () => {
  emit("close");
};

const handleOverlayClick = (e) => {
  if (e.target === e.currentTarget) {
    closePanel();
  }
};

watch(
  () => props.isOpen,
  (newVal) => {
    if (newVal) {
      fetchSettings();
    }
  },
);

onMounted(() => {
  if (props.isOpen) {
    fetchSettings();
  }
});
</script>

<template>
  <Teleport to="body">
    <Transition name="fade">
      <div v-if="isOpen" class="settings-overlay" @click="handleOverlayClick">
        <Transition name="slide">
          <div v-if="isOpen" class="settings-panel">
            <div class="settings-header">
              <h2>Settings</h2>
              <button class="close-btn" @click="closePanel" title="Close">
                <svg
                  width="20"
                  height="20"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                >
                  <line x1="18" y1="6" x2="6" y2="18" />
                  <line x1="6" y1="6" x2="18" y2="18" />
                </svg>
              </button>
            </div>

            <div class="settings-content">
              <div v-if="error" class="alert alert-error">
                {{ error }}
              </div>

              <div v-if="success" class="alert alert-success">
                {{ success }}
              </div>

              <div v-if="loading" class="loading-state">
                <div class="spinner"></div>
                <span>Loading settings...</span>
              </div>

              <div v-else class="settings-section">
                <h3 class="section-title">
                  <svg
                    width="18"
                    height="18"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                  >
                    <path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5" />
                  </svg>
                  MCP Integration
                </h3>
                <p class="section-description">
                  Model Context Protocol (MCP) allows AI assistants like Claude to interact with your DMARC data.
                  When enabled, MCP tools are available at <code>{{ settings.mcp_path }}</code>
                </p>

                <div class="setting-item">
                  <div class="setting-info">
                    <label class="setting-label">Enable MCP Server</label>
                    <span class="setting-hint">
                      Allow AI assistants to query your DMARC reports and statistics
                    </span>
                  </div>
                  <button
                    class="toggle-btn"
                    :class="{ active: settings.mcp_enabled }"
                    :disabled="saving"
                    @click="toggleMCP"
                  >
                    <span class="toggle-slider"></span>
                    <span class="toggle-label">{{ settings.mcp_enabled ? 'Enabled' : 'Disabled' }}</span>
                  </button>
                </div>

                <div v-if="settings.mcp_enabled" class="mcp-info">
                  <div class="info-box">
                    <svg
                      width="16"
                      height="16"
                      viewBox="0 0 24 24"
                      fill="none"
                      stroke="currentColor"
                      stroke-width="2"
                    >
                      <circle cx="12" cy="12" r="10" />
                      <line x1="12" y1="16" x2="12" y2="12" />
                      <line x1="12" y1="8" x2="12.01" y2="8" />
                    </svg>
                    <div class="info-content">
                      <strong>MCP Endpoint Active</strong>
                      <p>Configure your MCP client to connect to: <code>{{ window.location.origin }}{{ settings.mcp_path }}</code></p>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </Transition>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.settings-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(4px);
  z-index: 100;
  display: flex;
  justify-content: flex-end;
}

.settings-panel {
  width: 100%;
  max-width: 480px;
  background: var(--bg-card);
  height: 100%;
  box-shadow: -4px 0 24px rgba(0, 0, 0, 0.2);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.settings-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  border-bottom: 1px solid var(--border-subtle);
}

.settings-header h2 {
  margin: 0;
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--text-main);
}

.close-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border: none;
  background: transparent;
  color: var(--text-muted);
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s;
}

.close-btn:hover {
  background: var(--bg-app);
  color: var(--text-main);
}

.settings-content {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  padding: 40px;
  color: var(--text-muted);
}

.spinner {
  width: 32px;
  height: 32px;
  border: 3px solid var(--border-subtle);
  border-top-color: var(--c-primary);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.alert {
  padding: 12px 16px;
  border-radius: 8px;
  margin-bottom: 16px;
  font-size: 0.875rem;
}

.alert-error {
  background: var(--c-danger-soft);
  color: var(--c-danger);
  border: 1px solid var(--c-danger);
}

.alert-success {
  background: var(--c-success-soft);
  color: var(--c-success);
  border: 1px solid var(--c-success);
}

.settings-section {
  margin-bottom: 24px;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 1rem;
  font-weight: 600;
  color: var(--text-main);
  margin: 0 0 8px 0;
}

.section-title svg {
  color: var(--c-primary);
}

.section-description {
  color: var(--text-muted);
  font-size: 0.875rem;
  line-height: 1.5;
  margin: 0 0 20px 0;
}

.section-description code {
  background: var(--bg-app);
  padding: 2px 6px;
  border-radius: 4px;
  font-family: var(--font-mono);
  font-size: 0.8125rem;
}

.setting-item {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
  padding: 16px;
  background: var(--bg-app);
  border: 1px solid var(--border-subtle);
  border-radius: 8px;
}

.setting-info {
  flex: 1;
}

.setting-label {
  display: block;
  font-weight: 500;
  color: var(--text-main);
  margin-bottom: 4px;
}

.setting-hint {
  font-size: 0.8125rem;
  color: var(--text-muted);
}

.toggle-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  background: var(--bg-card);
  border: 1px solid var(--border-subtle);
  border-radius: 20px;
  cursor: pointer;
  transition: all 0.2s;
  flex-shrink: 0;
}

.toggle-btn:hover:not(:disabled) {
  border-color: var(--c-primary);
}

.toggle-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.toggle-slider {
  width: 36px;
  height: 20px;
  background: var(--border-subtle);
  border-radius: 10px;
  position: relative;
  transition: background 0.2s;
}

.toggle-slider::after {
  content: "";
  position: absolute;
  width: 16px;
  height: 16px;
  background: white;
  border-radius: 50%;
  top: 2px;
  left: 2px;
  transition: transform 0.2s;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.2);
}

.toggle-btn.active .toggle-slider {
  background: var(--c-primary);
}

.toggle-btn.active .toggle-slider::after {
  transform: translateX(16px);
}

.toggle-label {
  font-size: 0.8125rem;
  font-weight: 500;
  color: var(--text-muted);
  min-width: 60px;
}

.toggle-btn.active .toggle-label {
  color: var(--c-primary);
}

.mcp-info {
  margin-top: 16px;
}

.info-box {
  display: flex;
  gap: 12px;
  padding: 16px;
  background: var(--c-primary-soft, rgba(59, 130, 246, 0.1));
  border: 1px solid var(--c-primary);
  border-radius: 8px;
}

.info-box svg {
  color: var(--c-primary);
  flex-shrink: 0;
  margin-top: 2px;
}

.info-content {
  flex: 1;
}

.info-content strong {
  display: block;
  color: var(--text-main);
  margin-bottom: 4px;
}

.info-content p {
  margin: 0;
  font-size: 0.8125rem;
  color: var(--text-muted);
}

.info-content code {
  background: var(--bg-card);
  padding: 2px 6px;
  border-radius: 4px;
  font-family: var(--font-mono);
  font-size: 0.75rem;
  word-break: break-all;
}

/* Transitions */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.slide-enter-active,
.slide-leave-active {
  transition: transform 0.3s ease;
}

.slide-enter-from,
.slide-leave-to {
  transform: translateX(100%);
}

/* Responsive */
@media (max-width: 480px) {
  .settings-panel {
    max-width: 100%;
  }

  .setting-item {
    flex-direction: column;
    gap: 12px;
  }

  .toggle-btn {
    align-self: flex-start;
  }
}
</style>
