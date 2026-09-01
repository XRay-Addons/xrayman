#!/bin/bash
set -e

# запускается и при "apt remove", и при "apt purge" - $1 будет "remove" или "purge"
if [ "$1" = "purge" ]; then
  # полная зачистка - решите сами, нужно ли удалять данные/юзера при purge
  # rm -rf /var/lib/xray-node
  # rm -rf /etc/xray-node
  # userdel xray-node 2>/dev/null || true
  :
fi

systemctl daemon-reload 2>/dev/null || true