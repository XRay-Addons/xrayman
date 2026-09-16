#!/bin/sh
set -e
 
SERVICE_USER="xray-node"
SERVICE_GROUP="xray-node"
 
# Важно: этот скрипт выполняется ДО того, как пакетный менеджер распакует
# файлы. Файл xray-node.env в contents помечен owner: xray-node/group: xray-node,
# поэтому группа и пользователь должны существовать уже на этом этапе,
# иначе dpkg/rpm не смогут выставить владельца и упадут с ошибкой.
 
if ! getent group "$SERVICE_GROUP" >/dev/null 2>&1; then
    groupadd --system "$SERVICE_GROUP"
fi
 
if ! getent passwd "$SERVICE_USER" >/dev/null 2>&1; then
    useradd --system --gid "$SERVICE_GROUP" --no-create-home \
        --shell /usr/sbin/nologin "$SERVICE_USER"
fi
 
exit 0
