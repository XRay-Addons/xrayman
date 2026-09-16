#!/bin/sh
set -e
 
# Для deb: $1="purge" передаёт сам dpkg (это его механизм, а не nfpm) —
# conffile /etc/xray-node/xray-node.env к этому моменту dpkg уже удалил
# самостоятельно, поэтому дублировать rm для него не нужно.
# Для rpm аналогичная логика приходит через $1 (макросы %postun),
# nfpm сам подставляет нужный shell-код под нужный формат.
 
if [ "$1" = "purge" ] || [ "$1" = "0" ]; then
    CONFIG_DIR="/etc/xray-node"
 
    if [ -d "$CONFIG_DIR" ]; then
        echo "Purging configuration directory: $CONFIG_DIR"
        rm -rf "$CONFIG_DIR"
    fi
 
    if command -v systemctl >/dev/null 2>&1; then
        systemctl daemon-reload || true
    fi
fi
 
exit 0
 