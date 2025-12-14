#!/bin/bash
# Parse DMARC - DigitalOcean Marketplace Image Check
# This script validates that the image meets DigitalOcean Marketplace requirements
set -euo pipefail

echo "==> Running DigitalOcean Marketplace image checks..."

ERRORS=0

# Function to check and report
check() {
    local description="$1"
    local command="$2"

    if eval "$command" >/dev/null 2>&1; then
        echo "[PASS] $description"
    else
        echo "[FAIL] $description"
        ERRORS=$((ERRORS + 1))
    fi
}

echo ""
echo "=== Application Checks ==="
check "Parse DMARC binary exists" "[ -x /opt/parse-dmarc/parse-dmarc ]"
check "Parse DMARC symlink exists" "[ -L /usr/local/bin/parse-dmarc ]"
check "Parse DMARC systemd service exists" "[ -f /etc/systemd/system/parse-dmarc.service ]"
check "Parse DMARC user exists" "id parse-dmarc"
check "Parse DMARC data directory exists" "[ -d /var/lib/parse-dmarc ]"
check "Parse DMARC config directory exists" "[ -d /etc/parse-dmarc ]"
check "First boot script exists" "[ -x /opt/parse-dmarc/first-boot.sh ]"
check "MOTD script exists" "[ -x /etc/update-motd.d/99-parse-dmarc ]"

echo ""
echo "=== Security Checks ==="
check "UFW is enabled" "ufw status | grep -q 'Status: active'"
check "SSH is allowed" "ufw status | grep -q '22'"
check "Port 8080 is allowed" "ufw status | grep -q '8080'"
check "Fail2ban is installed" "systemctl is-enabled fail2ban"
check "Unattended upgrades configured" "[ -f /etc/apt/apt.conf.d/20auto-upgrades ]"

echo ""
echo "=== System Checks ==="
check "SSH host keys removed" "[ ! -f /etc/ssh/ssh_host_rsa_key ]"
check "Machine ID is empty" "[ ! -s /etc/machine-id ]"
check "Root bash history cleared" "[ ! -s /root/.bash_history ]"
check "Systemd service enabled" "systemctl is-enabled parse-dmarc"

echo ""
echo "=== Service Validation ==="
# Quick syntax check of the config
if /opt/parse-dmarc/parse-dmarc --help >/dev/null 2>&1; then
    echo "[PASS] Parse DMARC binary runs correctly"
else
    echo "[FAIL] Parse DMARC binary does not run"
    ERRORS=$((ERRORS + 1))
fi

echo ""
echo "==================================="
if [ $ERRORS -eq 0 ]; then
    echo "All checks passed! Image is ready for Marketplace."
    exit 0
else
    echo "WARNING: $ERRORS check(s) failed."
    echo "Please review and fix the issues before submitting to Marketplace."
    exit 1
fi
