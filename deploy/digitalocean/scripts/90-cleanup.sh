#!/bin/bash
# Parse DMARC - DigitalOcean Droplet Setup
# Script 90: Cleanup for image creation
set -euo pipefail

echo "==> Cleaning up for image creation..."

# Remove SSH keys (will be regenerated on first boot)
rm -f /etc/ssh/ssh_host_*

# Clean apt cache
apt-get -y autoremove
apt-get -y autoclean
apt-get -y clean

# Clear logs
find /var/log -type f -exec truncate -s 0 {} \;
journalctl --rotate
journalctl --vacuum-time=1s

# Clear temp files
rm -rf /tmp/*
rm -rf /var/tmp/*

# Clear bash history
history -c
cat /dev/null > ~/.bash_history
unset HISTFILE

# Remove machine-id (will be regenerated on first boot)
truncate -s 0 /etc/machine-id
rm -f /var/lib/dbus/machine-id

# Clear cloud-init state
# Only clean cloud-init if present; fail if cleaning fails
if command -v cloud-init >/dev/null 2>&1; then
    cloud-init clean --logs
fi

# Remove any sensitive data
rm -f /root/.ssh/authorized_keys
rm -f /root/.ssh/known_hosts

# Zero out free space for better compression (optional, can take time)
# dd if=/dev/zero of=/EMPTY bs=1M 2>/dev/null || true
# rm -f /EMPTY

echo "==> Cleanup complete"
