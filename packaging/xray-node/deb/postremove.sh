#!/bin/sh
set -e

# Для deb: $1="purge" передаёт сам dpkg (это его механизм, а не nfpm) —
# conffile /etc/xray-node/xray-node.env к этому моменту dpkg уже удалил
# самостоятельно, поэтому дублировать rm для него не нужно.
# Для rpm аналогичная логика приходит через $1 (макросы %postun),
# nfpm сам подставляет нужный shell-код под нужный формат.
#
# remove (apt remove)  — оставляет конфиги и persistent-данные (ключи и т.п.)
# purge  (apt purge)   — удаляет вообще всё, включая persistent-данные

# daemon-reload нужен при ЛЮБОМ удалении: файл .service к этому моменту
# dpkg уже стёр с диска сам (он не conffile и не наши данные — обычный
# файл пакета), но systemd продолжит "помнить" юнит, пока не перечитает
# конфигурацию.
if command -v systemctl >/dev/null 2>&1; then
    systemctl daemon-reload || true
fi

if [ "$1" = "purge" ] || [ "$1" = "0" ]; then
    CONFIG_DIR="/etc/xray-node"
    STATE_DIR="/var/lib/xray-node"

    if [ -d "$CONFIG_DIR" ]; then
        echo "Purging configuration directory: $CONFIG_DIR"
        rm -rf "$CONFIG_DIR"
    fi

    if [ -d "$STATE_DIR" ]; then
        echo "Purging state directory (incl. persistent keys): $STATE_DIR"
        rm -rf "$STATE_DIR"
    fi

    # Пользователя и группу тоже убираем при полном purge —
    # они больше никому не нужны, если пакет полностью удалён
    if getent passwd xray-node >/dev/null 2>&1; then
        userdel xray-node 2>/dev/null || true
    fi
    if getent group xray-node >/dev/null 2>&1; then
        groupdel xray-node 2>/dev/null || true
    fi
fi

exit 0