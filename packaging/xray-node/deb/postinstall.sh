#!/bin/sh
set -e

. /usr/lib/xray-node/utils.sh

SERVICE="xray-node"
SERVICE_USER="xray-node"
SERVICE_GROUP="xray-node"

echo "Run $SERVICE package postinstall.sh script..."
installSystemdService "$SERVICE"
exit 0
