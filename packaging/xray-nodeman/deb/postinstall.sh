#!/bin/sh
set -e

SERVICE_USER="xray-nodeman"
SERVICE_GROUP="xray-nodeman"

CONFIG_DIR="/etc/xray-nodeman"
ENV_FILE="${CONFIG_DIR}/xray-nodeman.env"
ENV_EXAMPLE="${CONFIG_DIR}/xray-nodeman.env.example"

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
# АВТОМАТИЧЕСКОЕ СОЗДАНИЕ КОНФИГА
CONFIG_CREATED=false

if [ ! -f "$ENV_FILE" ]; then
    if [ -f "$ENV_EXAMPLE" ]; then
        echo "Creating configuration file from example..."
        cp "$ENV_EXAMPLE" "$ENV_FILE"
        chown root:root "$ENV_FILE"
        chmod 600 "$ENV_FILE"
        CONFIG_CREATED=true
    else
        echo "Warning: Example file not found at $ENV_EXAMPLE. Cannot create config."
    fi
fi

echo "Reload systemd units."
if command -v systemctl >/dev/null 2>&1; then
    systemctl daemon-reload || true
fi

echo
echo "╔═════════════════════════════════════════════════════════════╗"
echo "║     XRayMan Node Manager installed successfully!            ║"
echo "╠═════════════════════════════════════════════════════════════╣"
echo "║ [!] Setup app configuration [!]                             ║"
echo "║   [!] User manual available via sudo apt show xray-nodeman  ║"
echo "║   [!] User manual available via sudo apt show xray-nodeman  ║"
echo "║   [!] User manual available via sudo apt show xray-nodeman  ║"
echo "║                                                             ║"
echo "║ Then start (or restart the service:                         ║"
echo "║   sudo systemctl enable --now xray-nodeman                  ║"
echo "║   or: sudo systemctl restart xray-nodeman                   ║"
echo "║                                                             ║"
echo "╚═════════════════════════════════════════════════════════════╝"

exit 0
