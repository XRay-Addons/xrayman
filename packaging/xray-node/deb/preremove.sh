#!/bin/sh
set -e

. /usr/lib/xray-node/utils.sh

SERVICE="xray-node"

if [ "$1" = "remove" ]; then
    echo "Run $SERVICE package preremove.sh script..."
    removeSystemdService $SERVICE
fi
exit 0