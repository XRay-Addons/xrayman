#!/bin/sh
set -e

. /usr/lib/xray-node/utils.sh

SERVICE="xray-node"
SERVICE_USER="xray-node"
SERVICE_GROUP="xray-node"

if [ "$1" = "purge" ]; then
    echo "Run $SERVICE package postremove.sh script..."
    removeServiceUser "$SERVICE_USER" "$SERVICE_GROUP"
fi

exit 0