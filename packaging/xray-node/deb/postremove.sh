#!/bin/sh
set -e

SERVICE="xray-node"
SERVICE_USER="xray-node"
SERVICE_GROUP="xray-node"

removeServiceUser() {
    local SERVICE_USER="$1"
    local SERVICE_GROUP="$2"

    if getent passwd "$SERVICE_USER" >/dev/null 2>&1; then
        echo "Removing system user '$SERVICE_USER'"
        userdel "$SERVICE_USER"
    else
        echo "System user '$SERVICE_USER' does not exist, skipping"
    fi

    if getent group "$SERVICE_GROUP" >/dev/null 2>&1; then
        echo "Removing system group '$SERVICE_GROUP'"
        groupdel "$SERVICE_GROUP"
    else
        echo "System group '$SERVICE_GROUP' does not exist, skipping"
    fi
}

if [ "$1" = "purge" ]; then
    echo "Run $SERVICE package postremove.sh script..."
    removeServiceUser "$SERVICE_USER" "$SERVICE_GROUP"
fi

exit 0