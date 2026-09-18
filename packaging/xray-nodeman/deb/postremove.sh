#!/bin/sh
set -e

SERVICE="xray-nodeman"
SERVICE_USER="xray-nodeman"
SERVICE_GROUP="xray-nodeman"

if [ "$1" = "purge" ]; then
    echo "Run $SERVICE package postremove.sh script..."
    
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
esac

exit 0