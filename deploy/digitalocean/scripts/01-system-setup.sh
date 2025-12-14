#!/bin/bash
# Parse DMARC - DigitalOcean Droplet Setup
# Script 01: System setup and base packages
set -euo pipefail

echo "==> Updating system packages..."
apt-get update
apt-get -y upgrade

echo "==> Installing required packages..."
apt-get -y install \
    curl \
    wget \
    jq \
    ufw \
    fail2ban \
    unattended-upgrades \
    apt-listchanges

echo "==> Configuring automatic security updates..."
cat > /etc/apt/apt.conf.d/20auto-upgrades <<EOF
APT::Periodic::Update-Package-Lists "1";
APT::Periodic::Unattended-Upgrade "1";
APT::Periodic::AutocleanInterval "7";
EOF

echo "==> Configuring firewall..."
ufw default deny incoming
ufw default allow outgoing
ufw allow ssh
ufw allow 8080/tcp comment 'Parse DMARC'
ufw --force enable

echo "==> Configuring fail2ban..."
systemctl enable fail2ban
systemctl start fail2ban

echo "==> Creating parse-dmarc user..."
useradd --system --shell /usr/sbin/nologin --home-dir /var/lib/parse-dmarc parse-dmarc || true

echo "==> Creating directories..."
mkdir -p /var/lib/parse-dmarc
mkdir -p /etc/parse-dmarc
mkdir -p /opt/parse-dmarc

chown -R parse-dmarc:parse-dmarc /var/lib/parse-dmarc

echo "==> System setup complete"
