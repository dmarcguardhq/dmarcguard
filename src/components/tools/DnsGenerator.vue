<script setup>
import { reactive, computed, ref } from "vue";

// State for the DMARC record parts
const config = reactive({
  domain: "example.com",
  policy: "none",
  subdomainPolicy: "",
  emailRua: "",
  emailRuf: "",
  pct: 100,
  alignment: "r",
});

const copied = ref(false);
const copiedValue = ref(false);

// Generate the TXT record string
const recordValue = computed(() => {
  let parts = ["v=DMARC1"];

  // Policy
  parts.push(`p=${config.policy}`);

  // Subdomain Policy (only if different or explicitly set)
  if (config.subdomainPolicy && config.subdomainPolicy !== config.policy) {
    parts.push(`sp=${config.subdomainPolicy}`);
  }

  // Percentage
  if (config.pct < 100) {
    parts.push(`pct=${config.pct}`);
  }

  // Reporting URIs (handle mailto prefix logic)
  if (config.emailRua) {
    const emails = config.emailRua
      .split(",")
      .map((e) =>
        e.trim().startsWith("mailto:") ? e.trim() : `mailto:${e.trim()}`,
      );
    parts.push(`rua=${emails.join(",")}`);
  }

  if (config.emailRuf) {
    const emails = config.emailRuf
      .split(",")
      .map((e) =>
        e.trim().startsWith("mailto:") ? e.trim() : `mailto:${e.trim()}`,
      );
    parts.push(`ruf=${emails.join(",")}`);
  }

  // Alignment
  if (config.alignment === "s") {
    parts.push("adkim=s");
    parts.push("aspf=s");
  }

  return parts.join("; ");
});

// Full DNS record for copying
const fullRecord = computed(() => {
  return `_dmarc.${config.domain} IN TXT "${recordValue.value}"`;
});

function copyToClipboard() {
  navigator.clipboard
    .writeText(fullRecord.value)
    .then(function () {
      copied.value = true;
      setTimeout(function () {
        copied.value = false;
      }, 2000);
    })
    .catch(function (err) {
      console.error("Failed to copy:", err);
    });
}

function copyValueOnly() {
  navigator.clipboard
    .writeText(recordValue.value)
    .then(function () {
      copiedValue.value = true;
      setTimeout(function () {
        copiedValue.value = false;
      }, 2000);
    })
    .catch(function (err) {
      console.error("Failed to copy:", err);
    });
}
</script>

<template>
  <div class="generator-layout">
    <div class="card controls-card">
      <h2 class="card-title">
        <svg
          width="20"
          height="20"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
        >
          <rect x="3" y="3" width="18" height="18" rx="2" ry="2" />
          <line x1="3" y1="9" x2="21" y2="9" />
          <line x1="9" y1="21" x2="9" y2="9" />
        </svg>
        DMARC Record Generator
      </h2>

      <form class="config-form" @submit.prevent>
        <div class="form-group">
          <label for="domain">Target Domain</label>
          <input
            id="domain"
            type="text"
            v-model="config.domain"
            class="input-text"
            placeholder="example.com"
          />
          <p class="help-text">
            The domain you want to protect with DMARC. Reports will be sent to
            the email addresses you specify below (rua/ruf).
          </p>
        </div>

        <div class="form-group">
          <label>Enforcement Policy (p)</label>
          <div class="segmented-control">
            <button
              type="button"
              :class="{ active: config.policy === 'none' }"
              @click="config.policy = 'none'"
              title="Monitor only. No emails blocked."
            >
              None
            </button>
            <button
              type="button"
              :class="{ active: config.policy === 'quarantine' }"
              @click="config.policy = 'quarantine'"
              title="Send suspicious mail to spam folder."
            >
              Quarantine
            </button>
            <button
              type="button"
              :class="{ active: config.policy === 'reject' }"
              @click="config.policy = 'reject'"
              title="Block suspicious mail completely."
            >
              Reject
            </button>
          </div>
          <p class="help-text" v-if="config.policy === 'none'">
            Suggested for initial monitoring phase.
          </p>
          <p class="help-text warning-text" v-if="config.policy === 'reject'">
            Ensure your SPF/DKIM is perfect before rejecting.
          </p>
        </div>

        <div class="form-group">
          <label>Subdomain Policy (sp)</label>
          <div class="segmented-control">
            <button
              type="button"
              :class="{ active: !config.subdomainPolicy }"
              @click="config.subdomainPolicy = ''"
              title="Use same policy as main domain."
            >
              Same
            </button>
            <button
              type="button"
              :class="{ active: config.subdomainPolicy === 'none' }"
              @click="config.subdomainPolicy = 'none'"
            >
              None
            </button>
            <button
              type="button"
              :class="{ active: config.subdomainPolicy === 'quarantine' }"
              @click="config.subdomainPolicy = 'quarantine'"
            >
              Quarantine
            </button>
            <button
              type="button"
              :class="{ active: config.subdomainPolicy === 'reject' }"
              @click="config.subdomainPolicy = 'reject'"
            >
              Reject
            </button>
          </div>
        </div>

        <div class="form-group">
          <label for="rua">Aggregate Reports (rua)</label>
          <input
            id="rua"
            type="email"
            v-model="config.emailRua"
            class="input-text"
            placeholder="dmarc-reports@example.com"
          />
          <p class="help-text">Where Parse DMARC will fetch reports from.</p>
        </div>

        <div class="form-group">
          <label for="ruf">Failure Reports (ruf)</label>
          <input
            id="ruf"
            type="email"
            v-model="config.emailRuf"
            class="input-text"
            placeholder="dmarc-failures@example.com"
          />
          <p class="help-text">
            Optional. Receive immediate failure notifications.
          </p>
        </div>

        <div class="form-group">
          <label for="pct">Policy Percentage (pct)</label>
          <div class="range-wrapper">
            <input
              id="pct"
              type="range"
              v-model.number="config.pct"
              min="0"
              max="100"
              step="5"
              class="input-range"
            />
            <span class="range-value">{{ config.pct }}%</span>
          </div>
          <p class="help-text">Percentage of messages to apply policy to.</p>
        </div>

        <div class="form-group">
          <div class="check-group">
            <input
              type="checkbox"
              id="strict"
              :checked="config.alignment === 's'"
              @change="config.alignment = $event.target.checked ? 's' : 'r'"
            />
            <label for="strict">Strict Alignment (adkim=s, aspf=s)</label>
          </div>
          <p class="help-text">Require exact domain match for SPF and DKIM.</p>
        </div>
      </form>
    </div>

    <div class="card preview-card">
      <h2 class="card-title">
        <svg
          width="20"
          height="20"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
        >
          <polyline points="4 17 10 11 4 5" />
          <line x1="12" y1="19" x2="20" y2="19" />
        </svg>
        DNS Record Preview
      </h2>

      <div class="terminal-window">
        <div class="terminal-header">
          <span class="dot red"></span>
          <span class="dot yellow"></span>
          <span class="dot green"></span>
          <span class="terminal-title">TXT Record</span>
        </div>

        <div class="code-content">
          <div class="line">
            <span class="key">Host:</span>
            <span class="value">_dmarc</span>
          </div>
          <div class="line">
            <span class="key">Type:</span>
            <span class="value">TXT</span>
          </div>
          <div class="line break-word">
            <span class="key">Value:</span>
            <span class="string">"{{ recordValue }}"</span>
            <button
              class="btn-copy-inline"
              @click="copyValueOnly"
              :class="{ copied: copiedValue }"
              title="Copy value only"
            >
              <svg
                v-if="!copiedValue"
                width="14"
                height="14"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
              >
                <rect x="9" y="9" width="13" height="13" rx="2" ry="2" />
                <path
                  d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"
                />
              </svg>
              <svg
                v-else
                width="14"
                height="14"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
              >
                <polyline points="20 6 9 17 4 12" />
              </svg>
            </button>
          </div>
        </div>

        <button
          class="btn-copy"
          @click="copyToClipboard"
          :class="{ copied: copied }"
        >
          <svg
            v-if="!copied"
            width="16"
            height="16"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
          >
            <rect x="9" y="9" width="13" height="13" rx="2" ry="2" />
            <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1" />
          </svg>
          <svg
            v-else
            width="16"
            height="16"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
          >
            <polyline points="20 6 9 17 4 12" />
          </svg>
          {{ copied ? "Copied!" : "Copy Full Record" }}
        </button>
      </div>

      <div class="instructions">
        <h3>
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
          Next Steps
        </h3>
        <ol>
          <li>Log in to your DNS provider (Cloudflare, GoDaddy, etc).</li>
          <li>Create a new <strong>TXT</strong> record.</li>
          <li>
            Set name/host to <code>_dmarc</code> (or
            <code>_dmarc.{{ config.domain }}</code
            >).
          </li>
          <li>Paste the value from above.</li>
          <li>Wait for DNS propagation (up to 48 hours).</li>
        </ol>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* --- Layout --- */
.generator-layout {
  display: grid;
  grid-template-columns: 1fr;
  gap: 24px;
}

@media (min-width: 900px) {
  .generator-layout {
    grid-template-columns: 1fr 1fr;
    align-items: start;
  }
}

.card {
  background: var(--bg-card);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-card);
  padding: 24px;
}

.card-title {
  margin: 0 0 24px 0;
  font-size: 1rem;
  font-weight: 600;
  color: var(--text-main);
  display: flex;
  align-items: center;
  gap: 10px;
}

.card-title svg {
  color: var(--c-primary);
}

/* --- Form Elements --- */
.config-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-group label {
  display: block;
  font-size: 0.875rem;
  font-weight: 500;
  margin-bottom: 8px;
  color: var(--text-main);
}

.input-text {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid var(--border-subtle);
  border-radius: 6px;
  background: var(--bg-app);
  color: var(--text-main);
  font-family: var(--font-mono);
  font-size: 0.875rem;
  transition: border-color 0.2s;
}

.input-text::placeholder {
  color: var(--text-muted);
}

.input-text:focus {
  outline: none;
  border-color: var(--c-primary);
}

/* Segmented Control */
.segmented-control {
  display: flex;
  background: var(--bg-app);
  padding: 4px;
  border-radius: 8px;
  border: 1px solid var(--border-subtle);
  gap: 4px;
}

.segmented-control button {
  flex: 1;
  background: transparent;
  border: none;
  padding: 8px 12px;
  font-size: 0.8125rem;
  font-weight: 500;
  color: var(--text-muted);
  cursor: pointer;
  border-radius: 6px;
  transition: all 0.2s;
  font-family: var(--font-sans);
}

.segmented-control button:hover {
  color: var(--text-main);
}

.segmented-control button.active {
  background: var(--bg-card);
  color: var(--c-primary);
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  font-weight: 600;
}

[data-theme="dark"] .segmented-control button.active {
  background: var(--c-primary);
  color: white;
}

.help-text {
  font-size: 0.75rem;
  color: var(--text-muted);
  margin: 6px 0 0 0;
}

.warning-text {
  color: var(--c-warning);
}

/* Range Slider */
.range-wrapper {
  display: flex;
  align-items: center;
  gap: 12px;
}

.input-range {
  flex: 1;
  height: 6px;
  -webkit-appearance: none;
  appearance: none;
  background: var(--border-subtle);
  border-radius: 3px;
  outline: none;
}

.input-range::-webkit-slider-thumb {
  -webkit-appearance: none;
  appearance: none;
  width: 18px;
  height: 18px;
  background: var(--c-primary);
  border-radius: 50%;
  cursor: pointer;
  transition: transform 0.2s;
}

.input-range::-webkit-slider-thumb:hover {
  transform: scale(1.1);
}

.input-range::-moz-range-thumb {
  width: 18px;
  height: 18px;
  background: var(--c-primary);
  border-radius: 50%;
  cursor: pointer;
  border: none;
}

.range-value {
  min-width: 48px;
  text-align: right;
  font-family: var(--font-mono);
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--c-primary);
}

/* Checkbox */
.check-group {
  display: flex;
  align-items: center;
  gap: 10px;
}

.check-group input[type="checkbox"] {
  width: 18px;
  height: 18px;
  accent-color: var(--c-primary);
  cursor: pointer;
}

.check-group label {
  margin: 0;
  cursor: pointer;
}

/* --- Terminal Preview --- */
.preview-card {
  display: flex;
  flex-direction: column;
}

.terminal-window {
  background: var(--terminal-bg);
  border: 1px solid var(--border-subtle);
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  margin-bottom: 24px;
}

.terminal-header {
  background: var(--terminal-header-bg);
  padding: 10px 16px;
  display: flex;
  align-items: center;
  gap: 6px;
  border-bottom: 1px solid var(--border-subtle);
}

.dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  opacity: 0.6;
}
.red {
  background: #ef4444;
}
.yellow {
  background: #f59e0b;
}
.green {
  background: #10b981;
}

.terminal-title {
  margin-left: 8px;
  font-size: 0.75rem;
  color: var(--text-muted);
  font-family: var(--font-sans);
  font-weight: 500;
}

.code-content {
  padding: 20px;
  font-family: var(--font-mono);
  font-size: 0.875rem;
  color: var(--terminal-text);
  background: var(--terminal-code-bg);
}

.line {
  margin-bottom: 8px;
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  align-items: center;
}

.key {
  color: var(--terminal-key);
  min-width: 50px;
  font-weight: 500;
}

.value {
  color: var(--terminal-value);
  font-weight: 500;
}

.string {
  color: var(--terminal-string);
  word-break: break-all;
}

.btn-copy {
  width: 100%;
  background: var(--terminal-btn-bg);
  border: none;
  color: var(--terminal-btn-text);
  padding: 14px;
  cursor: pointer;
  font-weight: 600;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  transition:
    background 0.2s,
    color 0.2s;
  border-top: 1px solid var(--border-subtle);
  font-family: var(--font-sans);
  font-size: 0.875rem;
}

.btn-copy:hover {
  background: var(--terminal-btn-hover);
}

.btn-copy.copied {
  background: var(--c-success);
  color: white;
}

.btn-copy-inline {
  background: var(--terminal-inline-btn-bg);
  border: 1px solid var(--terminal-inline-btn-border);
  border-radius: 4px;
  padding: 4px 6px;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
  color: var(--terminal-inline-btn-text);
  margin-left: 8px;
}

.btn-copy-inline:hover {
  background: var(--terminal-inline-btn-hover);
  border-color: var(--terminal-inline-btn-border-hover);
}

.btn-copy-inline.copied {
  background: var(--c-success);
  border-color: var(--c-success);
  color: white;
}

/* --- Instructions --- */
.instructions h3 {
  font-size: 0.875rem;
  color: var(--text-main);
  margin-bottom: 12px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.instructions h3 svg {
  color: var(--text-muted);
}

.instructions ol {
  padding-left: 20px;
  color: var(--text-muted);
  font-size: 0.875rem;
  line-height: 1.8;
  margin: 0;
}

.instructions code {
  background: var(--bg-app);
  padding: 2px 6px;
  border-radius: 4px;
  font-family: var(--font-mono);
  font-size: 0.8em;
  border: 1px solid var(--border-subtle);
  color: var(--text-main);
}

/* --- Responsive --- */
@media (max-width: 899px) {
  .preview-card {
    order: -1;
  }
}
</style>
