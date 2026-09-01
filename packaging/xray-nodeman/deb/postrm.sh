#!/bin/sh
set -e

# $1 — это первый аргумент, который apt/dpkg передает скрипту.
# Значение "purge" означает, что пользователь выполнил `apt purge` или `dpkg --purge`.
if [ "$1" = "purge" ]; then
    CONFIG_DIR="/etc/xray-nodeman"

    if [ -d "$CONFIG_DIR" ]; then
        echo "Purging configuration directory: $CONFIG_DIR"
        # Удаляем созданный вручную .env и саму папку
        rm -rf "$CONFIG_DIR"
    fi
fi

exit 0
