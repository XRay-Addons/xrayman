#!/bin/sh
set -e

SERVICE="xray-node"
SERVICE_USER="xray-node"
SERVICE_GROUP="xray-node"

removeNonPackageData() {
    local SERVICE="$1"
    local PERSISTENT="/var/lib/$SERVICE/persistent"
    local CONFIG="/etc/$SERVICE/config/"

    if [ -d "$PERSISTENT" ]; then
        echo "Removing $SERVICE persistent data from '$PERSISTENT'"
        find "$PERSISTENT" -mindepth 1 -delete # only content, not dir itself
    else
        echo "No persistent data dir '$PERSISTENT', nothing to remove"
    fi

    if [ -d "$CONFIG" ]; then
        echo "Removing $SERVICE config data from '$CONFIG'"
        find "$CONFIG" -mindepth 1 -delete # only content, not dir itself
    else
        echo "No config data dir '$CONFIG', nothing to remove"
    fi
}

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

removeServiceEnvFile() {
    local SERVICE="$1"
    local ENV_FILE="/etc/$SERVICE/$SERVICE.env"
    if [ -e "$ENV_FILE" ]; then
        rm "$ENV_FILE" 
        echo "$SERVICE env file '$ENV_FILE' removed"
    fi
}

if [ "$1" = "purge" ]; then
    echo "Run $SERVICE package postremove.sh script..."
    removeNonPackageData  "$SERVICE"
    removeServiceEnvFile "$SERVICE"
    removeServiceUser "$SERVICE_USER" "$SERVICE_GROUP"
fi

exit 0