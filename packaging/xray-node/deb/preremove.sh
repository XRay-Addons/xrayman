#!/bin/sh
set -e
 
if command -v systemctl >/dev/null 2>&1; then
    if systemctl is-active --quiet xray-node.service; then
        echo "Stopping xray-node service..."
        systemctl stop xray-node.service || true
    fi
 
    echo "Disabling xray-node service..."
    systemctl disable xray-node.service || true
fi
 
exit 0
 