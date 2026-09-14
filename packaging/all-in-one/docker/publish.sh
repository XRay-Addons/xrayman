#!/bin/bash
. ./packaging/common/utils.sh

set -e

CLEAN_VER="${1:-0.0.0}"

# docker compose with replaced version placeholder
replace \
    ./packaging/all-in-one/docker/docker-compose.yml \
    ./build/docker/all-in-one/docker-compose.yml \
    VERSION_PLACEHOLDER "$CLEAN_VER"

# config xray-node example folder
copy ./packaging/xray-node/config-example ./build/docker/all-in-one/xray-node-config/example

# copy prometheus config
copy ./packaging/all-in-one/docker/datasource.yml ./build/docker/all-in-one/datasource.yml
copy ./packaging/all-in-one/docker/prometheus.yml ./build/docker/all-in-one/prometheus.yml
       
# docker-compose archive
tar -czf docker-compose.all-in-one.tar.gz -C ./build/docker/all-in-one .
