#!/bin/sh
set -e

SERVICE_USER="xray-node"
SERVICE_GROUP="xray-node"

CONFIG_DIR="/etc/xray-node"
ENV_FILE="${CONFIG_DIR}/xray-node.env"
ENV_EXAMPLE="${CONFIG_DIR}/xray-node.env.example"

XRAY_CONFIG_DIR="/etc/xray-node"
XRAY_GEODATA_DIR="/var/lib/geodata"
XRAY_BINARY_DIR="/usr/bin"

echo "Create system group."
if ! getent group "$SERVICE_GROUP" >/dev/null 2>&1; then
    groupadd --system "$SERVICE_GROUP"
fi

echo "Create system user."
if ! getent passwd "$SERVICE_USER" >/dev/null 2>&1; then
    useradd --system --gid "$SERVICE_GROUP" --no-create-home --shell /usr/sbin/nologin "$SERVICE_USER"
fi

# Выставляем безопасные права на директорию
if [ -d "$CONFIG_DIR" ]; then
    chown root:root "$CONFIG_DIR"
    chmod 755 "$CONFIG_DIR"
fi

# Генерим ключи в примеры конфигов и удаляем файл-кейген, он больше не нужен
echo "Generate config example key."
XRAY_CONFIG_DIR/xray_keygen.tmp.sh \
    /opt/xray-node/bin/xray \
    /etc/xray-node/xray_server.example.json:/etc/xray-node/xray_client.example.json \
    X25519_PRIVATE_KEY X25519_PUBLIC_KEY
rm XRAY_CONFIG_DIR/xray_keygen.tmp.sh

echo "Reload systemd units."
if command -v systemctl >/dev/null 2>&1; then
    systemctl daemon-reload || true
fi

echo
echo "╔═════════════════════════════════════════════════════════════╗"
echo "║             XRayMan Node installed successfully!            ║"
echo "╠═════════════════════════════════════════════════════════════╣"
echo "║ [!] Setup app configuration [!]                             ║"
echo "║   [!] User manual available via sudo apt show xray-node     ║"
echo "║   [!] User manual available via sudo apt show xray-node     ║"
echo "║   [!] User manual available via sudo apt show xray-node.    ║"
echo "║                                                             ║"
echo "║ Then start (or restart) the service:                        ║"
echo "║   sudo systemctl enable --now xray-node                     ║"
echo "║   or: sudo systemctl restart xray-node.                     ║"
echo "║                                                             ║"
echo "╚═════════════════════════════════════════════════════════════╝"

exit 0
