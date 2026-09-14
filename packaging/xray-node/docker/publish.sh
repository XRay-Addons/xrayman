#!/bin/bash
. ./packaging/common/utils.sh

set -e

CLEAN_VER="${1:-0.0.0}"

# docker compose with replaced version placeholder
replace \
    ./packaging/xray-node/docker/docker-compose.yml \
    ./build/docker/xray-node/docker-compose.yml \
    VERSION_PLACEHOLDER "$CLEAN_VER"

# config example folder
copy ./packaging/xray-node/config-example ./build/docker/xray-node/xray-node-config/example

# docker-compose archive
tar -czf docker-compose.xray-node.tar.gz -C ./build/docker/xray-node .
