#!/bin/sh
set -e

# Проверяем, есть ли systemctl в системе, чтобы скрипт не упал
if command -v systemctl >/dev/null 2>&1; then
    # Проверяем, активен ли сервис, прежде чем его останавливать
    if systemctl is-active --quiet xray-nodeman.service; then
        echo "Stopping xray-nodeman service..."
        systemctl stop xray-nodeman.service || true
    fi
    
    echo "Disabling xray-nodeman service..."
    systemctl disable xray-nodeman.service || true
fi

exit 0
