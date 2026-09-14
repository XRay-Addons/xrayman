#!/bin/bash
. ./packaging/common/utils.sh

set -e

CLEAN_VER="${1:-0.0.0}"

# docker compose with replaced version placeholder
replace \
    ./packaging/xray-nodeman/docker/docker-compose.yml \
    ./build/docker/xray-nodeman/docker-compose.yml \
    VERSION_PLACEHOLDER "$CLEAN_VER"


# docker-compose archive
tar -czf docker-compose.xray-nodeman.tar.gz -C ./build/docker/xray-nodeman .
