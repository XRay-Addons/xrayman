#!/bin/sh
set -e

SERVICE="xray-node"

removeSystemdService() {
    local SERVICE="$1"

    if systemctl is-active --quiet $SERVICE.service; then
        echo "Stopping $SERVICE service"
        systemctl stop $SERVICE.service
    else
        echo "$SERVICE service is already stopped, skipping"
    fi

    if systemctl is-enabled --quiet "$SERVICE.service"; then
        echo "Disabling $SERVICE service"
        systemctl disable $SERVICE.service
    else
        echo "$SERVICE service is already disabled, skipping"
    fi
}

if [ "$1" = "remove" ]; then
    echo "Run $SERVICE package preremove.sh script..."
    removeSystemdService $SERVICE
fi
exit 0