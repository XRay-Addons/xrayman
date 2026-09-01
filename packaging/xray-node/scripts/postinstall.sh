#!/bin/bash
set -e

mkdir -p /var/lib/xray-node

chown -R root:xray-node /etc/xray-node /var/lib/xray-node
chmod 750 /etc/xray-node
chmod 770 /var/lib/xray-node

chown root:xray-node /etc/xray-node/xray-node.env
chmod 640 /etc/xray-node/xray-node.env

chown root:xray-node /etc/xray-node/xray_server.example.json
chmod 640 /etc/xray-node/xray_server.example.json
chown root:xray-node /etc/xray-node/xray_client.example.json
chmod 640 /etc/xray-node/xray_client.example.json

systemctl daemon-reload
systemctl enable xray-node

cat <<'EOF'

xray-node installed. Key paths:
  binary       /opt/xray-node/bin/xray-node
  data         /opt/xray-node/data
  config dir   /etc/xray-node
  persist dir  /var/lib/xray-node
  env file     /etc/xray-node/xray-node.env
  service unit /etc/systemd/system/xray-node.service

Edit env file and run:
  sudo systemctl restart xray-node
View logs:
  journalctl -u xray-node -f -n 50
EOF