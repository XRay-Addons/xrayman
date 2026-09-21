#!/bin/sh
set -e

. /usr/lib/xray-node/utils.sh

SERVICE="xray-node"
SERVICE_USER="xray-node"
SERVICE_GROUP="xray-node"

echo "Run $SERVICE package preinstall.sh script..."
createServiceUser $SERVICE_USER $SERVICE_GROUP
exit 0