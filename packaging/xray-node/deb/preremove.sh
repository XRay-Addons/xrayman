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

removeNonPackageData() {
    local SERVICE="$1"
    local PERSISTENT="/var/lib/$SERVICE/persistent"

    if [ -d "$DIR" ]; then
        echo "Removing $SERVICE persistent data from '$PERSISTENT'"
        find "$PERSISTENT" -mindepth 1 -delete # only content, not dir itself
    else
        echo "No persistent data dir '$PERSISTENT', nothing to remove"
    fi
}

if [ "$1" = "remove" ]; then
    echo "Run $SERVICE package preremove.sh script..."
    removeSystemdService $SERVICE
    removeNonPackageData $SERVICE
fi
exit 0