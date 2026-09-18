#!/bin/sh
set -e

SERVICE="xray-nodeman"
SERVICE_USER="xray-nodeman"
SERVICE_GROUP="xray-nodeman"

echo "Run $SERVICE package postinstall.sh script..."

DPKG_GARBAGE_FILE="/etc/$SERVICE/$SERVICE.env.dpkg-dist"
if [ -f "$DPKG_GARBAGE_FILE" ]; then
    echo "Removing dpkg garbage '$DPKG_GARBAGE_FILE'"
    rm -f "$DPKG_GARBAGE_FILE"
fi

if ! getent group $SERVICE_GROUP >/dev/null; then
    echo "Create system group '$SERVICE_GROUP'"
    groupadd --system "$SERVICE_GROUP"
else:
     echo "System group '$SERVICE_GROUP' already exists, skipping"
fi

echo "Create system user."
if ! getent passwd "$SERVICE_USER" >/dev/null 2>&1; then
    echo "Create system user '$SERVICE_USER'"
    useradd --system --gid "$SERVICE_GROUP" --no-create-home --shell /usr/sbin/nologin "$SERVICE_USER"
else:
    echo "System user '$SERVICE_USER' already exists, skipping"
fi

echo "Reload systemd units."
systemctl daemon-reload

echo "Enable $SERVICE service."
systemctl enable $SERVICE


printf '\n'
printf '╔════════════════════════════════════════════════════════════════╗\n'
printf '║     XRayMan %-16s service installed successfully!   ║\n' "$SERVICE"
printf '╠════════════════════════════════════════════════════════════════╣\n'
printf '║ [!] Setup app configuration [!]                                ║\n'
printf '║   [!] User manual available via sudo apt show %-16s ║\n' "$SERVICE"
printf '║                                                                ║\n'
printf '║ Then start/restart the service:                                ║\n'
printf '║   sudo systemctl restart %-16s                      ║\n' "$SERVICE"
printf '║                                                                ║\n'
printf '╚════════════════════════════════════════════════════════════════╝\n'
printf '\n'

exit 0
