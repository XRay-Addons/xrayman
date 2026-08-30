#!/bin/sh
set -e
SERVICE_USER="xray-nodeman"
SERVICE_GROUP="xray-nodeman"

CONFIG_DIR="/etc/xrayman"
ENV_FILE="${CONFIG_DIR}/xray-nodeman.env
ENV_EXAMPLE="${CONFIG_DIR}/xray-nodeman.env.example"

echo "Create system group."
if ! getent group "$SERVICE_GROUP" >/dev/null 2>&1; then
groupadd --system "$SERVICE_GROUP"
fi

echo "Create system user."
if ! getent passwd "$SERVICE_USER" >/dev/null 2>&1; then
useradd --system --gid "$SERVICE_GROUP" --no-create-home --shell /usr/sbin/nologin "$SERVICE_USER"
fi

echo "Reload systemd units."
if command -v systemctl >/dev/null 2>&1; then
systemctl daemon-reload || true
fi
echo "Installation instructions."
echo
echo "XRayMan Node Manager installed."
echo
echo "Configuration is required before starting the service."
echo
if [ ! -f "$ENV_FILE" ]; then
echo "Create the configuration file:"
echo
echo " sudo cp $ENV_EXAMPLE $ENV_FILE"
echo " sudo chmod 600 $ENV_FILE"
else
echo "Current environment file:"
echo "----------------------------------------"
cat "$ENV_FILE"
echo "----------------------------------------"
echo
echo "Example environment file:"
echo "----------------------------------------"
cat "$ENV_EXAMPLE"
echo "----------------------------------------"
fi

echo
echo "Then start the service:"
echo
echo " sudo systemctl enable --now xray-nodeman"
echo


exit 0