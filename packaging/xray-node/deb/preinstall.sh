#!/bin/sh
set -e

SERVICE="xray-node"
SERVICE_USER="xray-node"
SERVICE_GROUP="xray-node"


createServiceUser() {
    local SERVICE_USER="$1"
    local SERVICE_GROUP="$2"

    if ! getent group "$SERVICE_GROUP" >/dev/null; then
        echo "Create system group '$SERVICE_GROUP'"
        groupadd --system "$SERVICE_GROUP"
    else
        echo "System group '$SERVICE_GROUP' already exists, skipping"
    fi

    if ! getent passwd "$SERVICE_USER" >/dev/null 2>&1; then
        echo "Create system user '$SERVICE_USER'"
        useradd --system --gid "$SERVICE_GROUP" --no-create-home --shell /usr/sbin/nologin "$SERVICE_USER"
    else
        echo "System user '$SERVICE_USER' already exists, skipping"
    fi
}

echo "Run $SERVICE package preinstall.sh script..."
createServiceUser $SERVICE_USER $SERVICE_GROUP
exit 0