#!/bin/sh
set -e

SERVICE="xray-node"
SERVICE_USER="xray-node"
SERVICE_GROUP="xray-node"

installSystemdService() {
    local SERVICE="$1"

    DPKG_GARBAGE_FILE="/etc/$SERVICE/$SERVICE.env.dpkg-dist"
    if [ -f "$DPKG_GARBAGE_FILE" ]; then
        echo "Removing dpkg garbage '$DPKG_GARBAGE_FILE'"
        rm -f "$DPKG_GARBAGE_FILE"
    fi

    echo "Reload systemd units."
    systemctl daemon-reload

    echo "Enable $SERVICE service."
    systemctl enable $SERVICE


    printf '\n'
    printf '╔════════════════════════════════════════════════════════════════╗\n'
    printf '║     XRayMan service successfully installed: %-16s   ║\n' "$SERVICE"
    printf '╠════════════════════════════════════════════════════════════════╣\n'
    printf '║ [!] Setup app configuration [!]                                ║\n'
    printf '║   [!] User manual available via sudo apt show %-16s ║\n' "$SERVICE"
    printf '║                                                                ║\n'
    printf '║ Then start/restart the service:                                ║\n'
    printf '║   sudo systemctl restart %-16s                      ║\n' "$SERVICE"
    printf '║                                                                ║\n'
    printf '╚════════════════════════════════════════════════════════════════╝\n'
    printf '\n'
}

echo "Run $SERVICE package postinstall.sh script..."
installSystemdService "$SERVICE"
exit 0
