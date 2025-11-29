<template>
  <div class="dns-generator">
    <div class="generator-header">
      <h3 class="generator-title">
        <span class="generator-icon">🔧</span>
        DMARC Record Generator
      </h3>
      <p class="generator-subtitle">
        Generate the DNS TXT record for your domain's DMARC policy
      </p>
    </div>

    <div class="generator-content">
      <div class="form-section">
        <!-- Domain Input -->
        <div class="form-group">
          <label for="domain" class="form-label">
            Domain
            <span class="label-hint">Your domain name (e.g., example.com)</span>
          </label>
          <input
            id="domain"
            v-model="form.domain"
            type="text"
            class="form-input"
            placeholder="example.com"
            @input="validateDomain"
          />
          <span v-if="errors.domain" class="form-error">{{ errors.domain }}</span>
        </div>

        <!-- Policy Selection -->
        <div class="form-group">
          <label for="policy" class="form-label">
            Policy (p)
            <span class="label-hint">How receivers should handle messages that fail DMARC</span>
          </label>
          <select id="policy" v-model="form.policy" class="form-select">
            <option value="none">none - Monitor only, no action</option>
            <option value="quarantine">quarantine - Mark as spam</option>
            <option value="reject">reject - Block entirely</option>
          </select>
          <p class="form-help">
            <strong>Recommendation:</strong> Start with "none" to monitor, then gradually move to "quarantine" and "reject" after reviewing reports.
          </p>
        </div>

        <!-- Subdomain Policy -->
        <div class="form-group">
          <label for="sp" class="form-label">
            Subdomain Policy (sp)
            <span class="label-hint">Policy for subdomains (optional)</span>
          </label>
          <select id="sp" v-model="form.sp" class="form-select">
            <option value="">Inherit from main policy</option>
            <option value="none">none - Monitor only</option>
            <option value="quarantine">quarantine - Mark as spam</option>
            <option value="reject">reject - Block entirely</option>
          </select>
        </div>

        <!-- Reporting Address -->
        <div class="form-group">
          <label for="rua" class="form-label">
            Aggregate Report Address (rua)
            <span class="label-hint">Where to send aggregate reports</span>
          </label>
          <input
            id="rua"
            v-model="form.rua"
            type="email"
            class="form-input"
            placeholder="dmarc@example.com"
            @input="validateRua"
          />
          <span v-if="errors.rua" class="form-error">{{ errors.rua }}</span>
          <p class="form-help">
            This is where email providers will send your DMARC reports. Use the email address that Parse DMARC monitors.
          </p>
        </div>

        <!-- Advanced Options Toggle -->
        <button
          type="button"
          class="advanced-toggle"
          @click="showAdvanced = !showAdvanced"
          :aria-expanded="showAdvanced"
          aria-controls="advanced-options-content"
        >
          <span class="toggle-icon">{{ showAdvanced ? '▼' : '▶' }}</span>
          Advanced Options
        </button>

        <div v-if="showAdvanced" class="advanced-options" id="advanced-options-content">
          <!-- Percentage -->
          <div class="form-group">
            <label for="pct" class="form-label">
              Percentage (pct)
              <span class="label-hint">Percentage of messages to apply policy to</span>
            </label>
            <div class="range-input">
              <input
                id="pct"
                v-model.number="form.pct"
                type="range"
                min="1"
                max="100"
                class="form-range"
              />
              <span class="range-value">{{ form.pct }}%</span>
            </div>
            <p class="form-help">
              Use this for gradual rollout. Start at 10% when moving from "none" to "quarantine".
            </p>
          </div>

          <!-- DKIM Alignment -->
          <div class="form-group">
            <label for="adkim" class="form-label">
              DKIM Alignment (adkim)
              <span class="label-hint">How strictly to match DKIM domains</span>
            </label>
            <select id="adkim" v-model="form.adkim" class="form-select">
              <option value="r">relaxed - Allow subdomains</option>
              <option value="s">strict - Exact match only</option>
            </select>
          </div>

          <!-- SPF Alignment -->
          <div class="form-group">
            <label for="aspf" class="form-label">
              SPF Alignment (aspf)
              <span class="label-hint">How strictly to match SPF domains</span>
            </label>
            <select id="aspf" v-model="form.aspf" class="form-select">
              <option value="r">relaxed - Allow subdomains</option>
              <option value="s">strict - Exact match only</option>
            </select>
          </div>

          <!-- Forensic Reports -->
          <div class="form-group">
            <label for="ruf" class="form-label">
              Forensic Report Address (ruf)
              <span class="label-hint">Where to send failure reports (optional)</span>
            </label>
            <input
              id="ruf"
              v-model="form.ruf"
              type="email"
              class="form-input"
              placeholder="dmarc-forensic@example.com"
              @input="validateRuf"
            />
            <span v-if="errors.ruf" class="form-error">{{ errors.ruf }}</span>
            <p class="form-help">
              Forensic reports contain details about individual failures. Many providers don't send these due to privacy concerns.
            </p>
          </div>

          <!-- Failure Reporting Options -->
          <div class="form-group">
            <label for="fo" class="form-label">
              Failure Reporting (fo)
              <span class="label-hint">When to generate failure reports</span>
            </label>
            <select id="fo" v-model="form.fo" class="form-select">
              <option value="0">0 - Report if all checks fail (default)</option>
              <option value="1">1 - Report if any check fails</option>
              <option value="d">d - Report DKIM failures</option>
              <option value="s">s - Report SPF failures</option>
            </select>
          </div>

          <!-- Report Interval -->
          <div class="form-group">
            <label for="ri" class="form-label">
              Report Interval (ri)
              <span class="label-hint">How often to receive reports</span>
            </label>
            <select id="ri" v-model="form.ri" class="form-select">
              <option value="86400">Daily (86400 seconds)</option>
              <option value="43200">Twice daily (43200 seconds)</option>
              <option value="3600">Hourly (3600 seconds)</option>
            </select>
            <p class="form-help">
              Most providers send reports daily regardless of this setting.
            </p>
          </div>
        </div>
      </div>

      <!-- Preview Section -->
      <div class="preview-section">
        <div class="preview-header">
          <h4 class="preview-title">Generated Record</h4>
          <button
            class="copy-button"
            @click="copyRecord"
            :disabled="!isValid"
            :class="{ copied: copySuccess }"
          >
            <span v-if="copySuccess">✓ Copied!</span>
            <span v-else>📋 Copy</span>
          </button>
        </div>

        <div class="preview-content">
          <div class="dns-record-box">
            <div class="record-row">
              <span class="record-label">Name:</span>
              <code class="record-value">_dmarc.{{ form.domain || 'example.com' }}</code>
            </div>
            <div class="record-row">
              <span class="record-label">Type:</span>
              <code class="record-value">TXT</code>
            </div>
            <div class="record-row">
              <span class="record-label">Value:</span>
              <code class="record-value record-txt">{{ generatedRecord }}</code>
            </div>
          </div>
        </div>

        <!-- Provider Instructions -->
        <div class="provider-section">
          <button
            type="button"
            class="provider-toggle"
            @click="showProviders = !showProviders"
            :aria-expanded="showProviders"
            aria-controls="provider-list-content"
          >
            <span class="toggle-icon">{{ showProviders ? '▼' : '▶' }}</span>
            Provider-Specific Instructions
          </button>

          <div v-if="showProviders" class="provider-list" id="provider-list-content">
            <div class="provider-item">
              <h5 class="provider-name">Cloudflare</h5>
              <ol class="provider-steps">
                <li>Go to your domain's DNS settings</li>
                <li>Click "Add record"</li>
                <li>Type: <code>TXT</code></li>
                <li>Name: <code>_dmarc</code></li>
                <li>Content: paste the Value above</li>
                <li>Click "Save"</li>
              </ol>
            </div>

            <div class="provider-item">
              <h5 class="provider-name">AWS Route 53</h5>
              <ol class="provider-steps">
                <li>Go to your hosted zone</li>
                <li>Click "Create record"</li>
                <li>Record name: <code>_dmarc</code></li>
                <li>Record type: <code>TXT</code></li>
                <li>Value: paste the Value above (include quotes)</li>
                <li>Click "Create records"</li>
              </ol>
            </div>

            <div class="provider-item">
              <h5 class="provider-name">Google Domains / Squarespace</h5>
              <ol class="provider-steps">
                <li>Go to DNS settings</li>
                <li>Click "Manage custom records" or "Add"</li>
                <li>Host name: <code>_dmarc</code></li>
                <li>Type: <code>TXT</code></li>
                <li>Data: paste the Value above</li>
                <li>Save</li>
              </ol>
            </div>

            <div class="provider-item">
              <h5 class="provider-name">GoDaddy</h5>
              <ol class="provider-steps">
                <li>Go to "My Products" → Domain → "DNS"</li>
                <li>Click "Add" in the Records section</li>
                <li>Type: <code>TXT</code></li>
                <li>Name: <code>_dmarc</code></li>
                <li>Value: paste the Value above</li>
                <li>TTL: 1 Hour</li>
                <li>Click "Save"</li>
              </ol>
            </div>

            <div class="provider-item">
              <h5 class="provider-name">Namecheap</h5>
              <ol class="provider-steps">
                <li>Go to Domain List → Manage → Advanced DNS</li>
                <li>Click "Add New Record"</li>
                <li>Type: <code>TXT Record</code></li>
                <li>Host: <code>_dmarc</code></li>
                <li>Value: paste the Value above</li>
                <li>Click the checkmark to save</li>
              </ol>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { computed, reactive, ref } from "vue";

export default {
  name: "DnsGenerator",
  setup() {
    var form = reactive({
      domain: "",
      policy: "none",
      sp: "",
      rua: "",
      ruf: "",
      pct: 100,
      adkim: "r",
      aspf: "r",
      fo: "0",
      ri: "86400",
    });

    var errors = reactive({
      domain: "",
      rua: "",
      ruf: "",
    });

    var showAdvanced = ref(false);
    var showProviders = ref(false);
    var copySuccess = ref(false);

    function isValidEmail(email) {
      if (!email) return true; // Empty is valid (optional field)
      var emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
      return emailRegex.test(email);
    }

    function isValidDomain(domain) {
      if (!domain) return true; // Empty handled separately
      var domainRegex = /^[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z]{2,})+$/;
      return domainRegex.test(domain);
    }

    function validateDomain() {
      if (!form.domain) {
        errors.domain = "";
      } else if (!isValidDomain(form.domain)) {
        errors.domain = "Please enter a valid domain name";
      } else {
        errors.domain = "";
      }
    }

    function validateRua() {
      if (!form.rua) {
        errors.rua = "";
      } else if (!isValidEmail(form.rua)) {
        errors.rua = "Please enter a valid email address";
      } else {
        errors.rua = "";
      }
    }

    function validateRuf() {
      if (!form.ruf) {
        errors.ruf = "";
      } else if (!isValidEmail(form.ruf)) {
        errors.ruf = "Please enter a valid email address";
      } else {
        errors.ruf = "";
      }
    }

    var isValid = computed(function () {
      return (
        form.domain &&
        !errors.domain &&
        !errors.rua &&
        !errors.ruf
      );
    });

    var generatedRecord = computed(function () {
      var parts = ["v=DMARC1"];

      // Policy (required)
      parts.push("p=" + form.policy);

      // Subdomain policy (optional)
      if (form.sp) {
        parts.push("sp=" + form.sp);
      }

      // Aggregate reports (if provided)
      if (form.rua && !errors.rua) {
        parts.push("rua=mailto:" + form.rua);
      }

      // Forensic reports (if provided)
      if (form.ruf && !errors.ruf) {
        parts.push("ruf=mailto:" + form.ruf);
      }

      // Percentage (only if not 100)
      if (form.pct !== 100) {
        parts.push("pct=" + form.pct);
      }

      // DKIM alignment (only if strict)
      if (form.adkim === "s") {
        parts.push("adkim=s");
      }

      // SPF alignment (only if strict)
      if (form.aspf === "s") {
        parts.push("aspf=s");
      }

      // Failure options (only if not default)
      if (form.fo !== "0") {
        parts.push("fo=" + form.fo);
      }

      // Report interval (only if not daily)
      if (form.ri !== "86400") {
        parts.push("ri=" + form.ri);
      }

      return parts.join("; ");
    });

    function copyRecord() {
      if (!isValid.value) return;

      var textToCopy = generatedRecord.value;

      if (navigator.clipboard && navigator.clipboard.writeText) {
        navigator.clipboard.writeText(textToCopy).then(function () {
          copySuccess.value = true;
          setTimeout(function () {
            copySuccess.value = false;
          }, 2000);
        }).catch(function (err) {
          console.error("Failed to copy:", err);
          fallbackCopy(textToCopy);
        });
      } else {
        fallbackCopy(textToCopy);
      }
    }

    function fallbackCopy(text) {
      var textArea = document.createElement("textarea");
      textArea.value = text;
      textArea.style.position = "fixed";
      textArea.style.left = "-9999px";
      document.body.appendChild(textArea);
      textArea.select();
      try {
        document.execCommand("copy");
        copySuccess.value = true;
        setTimeout(function () {
          copySuccess.value = false;
        }, 2000);
      } catch (err) {
        console.error("Fallback copy failed:", err);
      }
      document.body.removeChild(textArea);
    }

    return {
      form,
      errors,
      showAdvanced,
      showProviders,
      copySuccess,
      isValid,
      generatedRecord,
      validateDomain,
      validateRua,
      validateRuf,
      copyRecord,
    };
  },
};
</script>

<style scoped>
.dns-generator {
  background: var(--color-bg-primary);
  border-radius: 12px;
  padding: 1.5rem;
  box-shadow: var(--shadow-sm);
  transition: var(--theme-transition);
}

.generator-header {
  margin-bottom: 1.5rem;
}

.generator-title {
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0 0 0.5rem 0;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.generator-icon {
  font-size: 1.25rem;
}

.generator-subtitle {
  color: var(--color-text-secondary);
  font-size: 0.9rem;
  margin: 0;
}

.generator-content {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 2rem;
}

@media (max-width: 900px) {
  .generator-content {
    grid-template-columns: 1fr;
  }
}

.form-section {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.form-label {
  font-weight: 600;
  color: var(--color-text-primary);
  font-size: 0.9rem;
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.label-hint {
  font-weight: 400;
  color: var(--color-text-secondary);
  font-size: 0.8rem;
}

.form-input,
.form-select {
  padding: 0.75rem;
  border: 1px solid var(--color-border);
  border-radius: 6px;
  font-size: 0.95rem;
  background: var(--color-bg-secondary);
  color: var(--color-text-primary);
  transition: border-color 0.2s, box-shadow 0.2s;
}

.form-input:focus,
.form-select:focus {
  outline: none;
  border-color: var(--color-accent);
  box-shadow: 0 0 0 3px var(--color-accent-light);
}

.form-input::placeholder {
  color: var(--color-text-muted);
}

.form-error {
  color: var(--color-error);
  font-size: 0.8rem;
}

.form-help {
  color: var(--color-text-secondary);
  font-size: 0.8rem;
  margin: 0;
  line-height: 1.4;
}

.form-help strong {
  color: var(--color-text-primary);
}

.advanced-toggle,
.provider-toggle {
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: 6px;
  padding: 0.75rem 1rem;
  font-size: 0.9rem;
  font-weight: 600;
  color: var(--color-text-primary);
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  transition: background 0.2s, border-color 0.2s;
  width: 100%;
  text-align: left;
}

.advanced-toggle:hover,
.provider-toggle:hover {
  background: var(--color-bg-tertiary);
  border-color: var(--color-accent);
}

.toggle-icon {
  font-size: 0.75rem;
  color: var(--color-text-secondary);
}

.advanced-options {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  padding: 1rem;
  background: var(--color-bg-secondary);
  border-radius: 6px;
  border: 1px solid var(--color-border);
}

.range-input {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.form-range {
  flex: 1;
  height: 6px;
  -webkit-appearance: none;
  appearance: none;
  background: var(--color-bg-tertiary);
  border-radius: 3px;
  cursor: pointer;
}

.form-range::-webkit-slider-thumb {
  -webkit-appearance: none;
  appearance: none;
  width: 18px;
  height: 18px;
  border-radius: 50%;
  background: var(--color-accent);
  cursor: pointer;
}

.form-range::-moz-range-thumb {
  width: 18px;
  height: 18px;
  border-radius: 50%;
  background: var(--color-accent);
  cursor: pointer;
  border: none;
}

.range-value {
  font-weight: 600;
  color: var(--color-accent);
  min-width: 3rem;
  text-align: right;
}

.preview-section {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.preview-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.preview-title {
  font-size: 1rem;
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0;
}

.copy-button {
  background: var(--color-accent);
  color: white;
  border: none;
  border-radius: 6px;
  padding: 0.5rem 1rem;
  font-size: 0.85rem;
  font-weight: 600;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  transition: background 0.2s, transform 0.1s;
}

.copy-button:hover:not(:disabled) {
  background: var(--color-accent-text);
  transform: translateY(-1px);
}

.copy-button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.copy-button.copied {
  background: var(--color-success);
}

.preview-content {
  background: var(--color-bg-secondary);
  border-radius: 6px;
  padding: 1rem;
  border: 1px solid var(--color-border);
}

.dns-record-box {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.record-row {
  display: flex;
  gap: 0.75rem;
  align-items: flex-start;
}

.record-label {
  font-weight: 600;
  color: var(--color-text-secondary);
  min-width: 50px;
  font-size: 0.85rem;
}

.record-value {
  background: var(--color-code-bg);
  color: var(--color-code-text);
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-family: "Courier New", monospace;
  font-size: 0.85rem;
  word-break: break-all;
}

.record-txt {
  flex: 1;
  line-height: 1.5;
}

.provider-section {
  margin-top: 0.5rem;
}

.provider-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  margin-top: 1rem;
  padding: 1rem;
  background: var(--color-bg-secondary);
  border-radius: 6px;
  border: 1px solid var(--color-border);
}

.provider-item {
  padding-bottom: 1rem;
  border-bottom: 1px solid var(--color-border);
}

.provider-item:last-child {
  padding-bottom: 0;
  border-bottom: none;
}

.provider-name {
  font-size: 0.95rem;
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0 0 0.5rem 0;
}

.provider-steps {
  margin: 0;
  padding-left: 1.25rem;
  color: var(--color-text-secondary);
  font-size: 0.85rem;
  line-height: 1.6;
}

.provider-steps code {
  background: var(--color-code-bg);
  color: var(--color-code-text);
  padding: 0.1rem 0.3rem;
  border-radius: 3px;
  font-family: "Courier New", monospace;
  font-size: 0.8rem;
}
</style>
