<script setup>
import { computed, ref, watch } from "vue";
import { useThemeStore, useSettingsStore } from "@/stores";

const props = defineProps({
  isOpen: {
    type: Boolean,
    default: false,
  },
});

const emit = defineEmits(["close"]);

const themeStore = useThemeStore();
const settingsStore = useSettingsStore();

// Local state for editing
const editingEndpoint = ref(settingsStore.apiEndpoint);
const validationError = ref("");

// Sync editing endpoint when modal opens or store changes
watch(
  () => props.isOpen,
  (open) => {
    if (open) {
      editingEndpoint.value = settingsStore.apiEndpoint;
      validationError.value = "";
    }
  },
);

// Validate on input
watch(editingEndpoint, (value) => {
  const result = settingsStore.validateApiEndpoint(value);
  validationError.value = result.valid ? "" : result.error;
});

// Computed states
const hasChanges = computed(() => {
  return editingEndpoint.value !== settingsStore.apiEndpoint;
});

const canSave = computed(() => {
  return hasChanges.value && !validationError.value;
});

const canTest = computed(() => {
  return !validationError.value && !settingsStore.isTestingConnection;
});

// Actions
function saveEndpoint() {
  if (!canSave.value) return;
  const result = settingsStore.setApiEndpoint(editingEndpoint.value);
  if (!result.success) {
    validationError.value = result.error;
  }
}

async function testEndpoint() {
  if (!canTest.value) return;
  await settingsStore.testConnection(editingEndpoint.value);
}

function resetEndpoint() {
  editingEndpoint.value = settingsStore.defaultApiEndpoint;
  settingsStore.resetApiEndpoint();
}

function closeModal() {
  emit("close");
}

function handleBackdropClick(event) {
  if (event.target === event.currentTarget) {
    closeModal();
  }
}

// Theme helpers
const themeOptions = [
  { value: "light", label: "Light", icon: "sun" },
  { value: "dark", label: "Dark", icon: "moon" },
  { value: "system", label: "System", icon: "monitor" },
];
</script>

<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="isOpen" class="modal-backdrop" @click="handleBackdropClick">
        <div
          class="modal"
          role="dialog"
          aria-modal="true"
          aria-labelledby="settings-title"
        >
          <div class="modal-header">
            <h2 id="settings-title" class="modal-title">
              <svg
                width="20"
                height="20"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
              >
                <circle cx="12" cy="12" r="3" />
                <path
                  d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"
                />
              </svg>
              Settings
            </h2>
            <button
              class="btn-close"
              @click="closeModal"
              aria-label="Close settings"
            >
              <svg
                width="20"
                height="20"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
              >
                <line x1="18" y1="6" x2="6" y2="18" />
                <line x1="6" y1="6" x2="18" y2="18" />
              </svg>
            </button>
          </div>

          <div class="modal-body">
            <!-- Theme Section -->
            <section class="settings-section">
              <h3 class="section-label">Theme</h3>
              <div class="theme-selector">
                <button
                  v-for="option in themeOptions"
                  :key="option.value"
                  class="theme-option"
                  :class="{ active: themeStore.theme === option.value }"
                  @click="themeStore.setTheme(option.value)"
                >
                  <!-- Sun icon -->
                  <svg
                    v-if="option.icon === 'sun'"
                    width="18"
                    height="18"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                  >
                    <circle cx="12" cy="12" r="4" />
                    <path d="M12 2v2" />
                    <path d="M12 20v2" />
                    <path d="m4.93 4.93 1.41 1.41" />
                    <path d="m17.66 17.66 1.41 1.41" />
                    <path d="M2 12h2" />
                    <path d="M20 12h2" />
                    <path d="m6.34 17.66-1.41 1.41" />
                    <path d="m19.07 4.93-1.41 1.41" />
                  </svg>
                  <!-- Moon icon -->
                  <svg
                    v-else-if="option.icon === 'moon'"
                    width="18"
                    height="18"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                  >
                    <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z" />
                  </svg>
                  <!-- Monitor icon -->
                  <svg
                    v-else
                    width="18"
                    height="18"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                  >
                    <rect x="2" y="3" width="20" height="14" rx="2" ry="2" />
                    <line x1="8" y1="21" x2="16" y2="21" />
                    <line x1="12" y1="17" x2="12" y2="21" />
                  </svg>
                  <span>{{ option.label }}</span>
                </button>
              </div>
            </section>

            <!-- API Endpoint Section -->
            <section class="settings-section">
              <h3 class="section-label">API Endpoint</h3>
              <p class="section-description">
                Configure the backend API URL. Changes take effect immediately.
              </p>

              <div class="input-group">
                <input
                  v-model="editingEndpoint"
                  type="url"
                  class="input"
                  :class="{ 'input-error': validationError }"
                  placeholder="https://example.com/api"
                  spellcheck="false"
                />
                <button
                  class="btn btn-test"
                  :disabled="!canTest"
                  @click="testEndpoint"
                >
                  <svg
                    v-if="settingsStore.isTestingConnection"
                    class="spinner"
                    width="16"
                    height="16"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                  >
                    <path d="M21 12a9 9 0 1 1-6.219-8.56" />
                  </svg>
                  <svg
                    v-else
                    width="16"
                    height="16"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                  >
                    <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14" />
                    <polyline points="22 4 12 14.01 9 11.01" />
                  </svg>
                  Test
                </button>
              </div>

              <!-- Validation error -->
              <p v-if="validationError" class="error-message">
                {{ validationError }}
              </p>

              <!-- Test result -->
              <div
                v-if="settingsStore.lastTestResult"
                class="test-result"
                :class="{
                  success: settingsStore.lastTestResult.success,
                  error: !settingsStore.lastTestResult.success,
                }"
              >
                <svg
                  v-if="settingsStore.lastTestResult.success"
                  width="16"
                  height="16"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                >
                  <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14" />
                  <polyline points="22 4 12 14.01 9 11.01" />
                </svg>
                <svg
                  v-else
                  width="16"
                  height="16"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                >
                  <circle cx="12" cy="12" r="10" />
                  <line x1="15" y1="9" x2="9" y2="15" />
                  <line x1="9" y1="9" x2="15" y2="15" />
                </svg>
                <span>{{ settingsStore.lastTestResult.message }}</span>
              </div>

              <!-- Actions -->
              <div class="endpoint-actions">
                <button
                  class="btn btn-secondary"
                  :disabled="!settingsStore.isCustomEndpoint"
                  @click="resetEndpoint"
                >
                  Reset to Default
                </button>
                <button
                  class="btn btn-primary"
                  :disabled="!canSave"
                  @click="saveEndpoint"
                >
                  Save Changes
                </button>
              </div>

              <p v-if="settingsStore.isCustomEndpoint" class="default-hint">
                Default: <code>{{ settingsStore.defaultApiEndpoint }}</code>
              </p>
            </section>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.modal-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 100;
  padding: 16px;
}

.modal {
  background: var(--bg-card);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-lg);
  width: 100%;
  max-width: 480px;
  max-height: 90vh;
  overflow: auto;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.2);
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid var(--border-subtle);
}

.modal-title {
  display: flex;
  align-items: center;
  gap: 10px;
  margin: 0;
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--text-main);
}

.modal-title svg {
  color: var(--text-muted);
}

.btn-close {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border: none;
  background: transparent;
  color: var(--text-muted);
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: all 0.2s;
}

.btn-close:hover {
  background: var(--bg-app);
  color: var(--text-main);
}

.modal-body {
  padding: 20px;
}

.settings-section {
  margin-bottom: 24px;
}

.settings-section:last-child {
  margin-bottom: 0;
}

.section-label {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--text-main);
  margin: 0 0 8px 0;
}

.section-description {
  font-size: 0.8125rem;
  color: var(--text-muted);
  margin: 0 0 12px 0;
  line-height: 1.5;
}

/* Theme Selector */
.theme-selector {
  display: flex;
  gap: 8px;
}

.theme-option {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  padding: 12px 8px;
  border: 1px solid var(--border-subtle);
  background: var(--bg-app);
  border-radius: var(--radius-card);
  color: var(--text-muted);
  font-size: 0.8125rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  font-family: var(--font-sans);
}

.theme-option:hover {
  border-color: var(--c-primary);
  color: var(--text-main);
}

.theme-option.active {
  background: var(--c-primary);
  border-color: var(--c-primary);
  color: white;
}

/* Input Group */
.input-group {
  display: flex;
  gap: 8px;
}

.input {
  flex: 1;
  padding: 10px 12px;
  border: 1px solid var(--border-subtle);
  background: var(--bg-app);
  border-radius: var(--radius-sm);
  color: var(--text-main);
  font-size: 0.875rem;
  font-family: var(--font-mono);
  transition: all 0.2s;
}

.input:focus {
  outline: none;
  border-color: var(--c-primary);
  box-shadow: 0 0 0 3px var(--c-primary-soft);
}

.input-error {
  border-color: var(--c-danger);
}

.input-error:focus {
  border-color: var(--c-danger);
  box-shadow: 0 0 0 3px rgba(239, 68, 68, 0.15);
}

/* Buttons */
.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 10px 16px;
  border: none;
  border-radius: var(--radius-sm);
  font-size: 0.875rem;
  font-weight: 500;
  font-family: var(--font-sans);
  cursor: pointer;
  transition: all 0.2s;
  white-space: nowrap;
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-test {
  background: var(--bg-app);
  border: 1px solid var(--border-subtle);
  color: var(--text-main);
}

.btn-test:hover:not(:disabled) {
  border-color: var(--c-primary);
  color: var(--c-primary);
}

.btn-primary {
  background: var(--c-primary);
  color: white;
}

.btn-primary:hover:not(:disabled) {
  filter: brightness(1.1);
}

.btn-secondary {
  background: var(--bg-app);
  border: 1px solid var(--border-subtle);
  color: var(--text-muted);
}

.btn-secondary:hover:not(:disabled) {
  border-color: var(--border-subtle);
  color: var(--text-main);
  background: var(--bg-card);
}

/* Error Message */
.error-message {
  margin: 8px 0 0 0;
  font-size: 0.8125rem;
  color: var(--c-danger);
}

/* Test Result */
.test-result {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 12px;
  padding: 10px 12px;
  border-radius: var(--radius-sm);
  font-size: 0.8125rem;
}

.test-result.success {
  background: rgba(16, 185, 129, 0.1);
  color: var(--c-success);
}

.test-result.error {
  background: rgba(239, 68, 68, 0.1);
  color: var(--c-danger);
}

/* Endpoint Actions */
.endpoint-actions {
  display: flex;
  gap: 8px;
  margin-top: 16px;
  justify-content: flex-end;
}

/* Default Hint */
.default-hint {
  margin: 12px 0 0 0;
  font-size: 0.75rem;
  color: var(--text-muted);
}

.default-hint code {
  font-family: var(--font-mono);
  background: var(--bg-app);
  padding: 2px 6px;
  border-radius: var(--radius-sm);
}

/* Spinner animation */
.spinner {
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

/* Modal transitions */
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.2s ease;
}

.modal-enter-active .modal,
.modal-leave-active .modal {
  transition:
    transform 0.2s ease,
    opacity 0.2s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-from .modal,
.modal-leave-to .modal {
  transform: scale(0.95);
  opacity: 0;
}

/* Responsive */
@media (max-width: 480px) {
  .modal {
    max-width: 100%;
    margin: 0;
    border-radius: var(--radius-card);
  }

  .theme-selector {
    flex-direction: column;
  }

  .input-group {
    flex-direction: column;
  }

  .endpoint-actions {
    flex-direction: column;
  }

  .btn {
    width: 100%;
  }
}
</style>
