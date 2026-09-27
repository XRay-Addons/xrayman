#!/bin/bash
set -e

copy() {
    mkdir -p "$(dirname "$2")"
    cp -r "$1" "$2"
    if [[ -n "$3" ]]; then
        chmod "$3" "$2"
    fi
}

replace() {
    local src="$1"
    local dst="$2"
    shift 2

    local args=()
    while (($#)); do
        args+=("-e" "s/$1/$2/")
        shift 2
    done

    mkdir -p "$(dirname "$dst")"
    sed "${args[@]}" "$src" > "$dst"
}

CLEAN_VER="${1:-0.0.0}"

# ╔═══════════════════════════════════╗ #
# ║                 node              ║ #
# ╚═══════════════════════════════════╝ #

# docker compose with replaced version placeholder
replace \
    ./packaging/docker-compose/xray-node/docker-compose.yml \
    ./build/docker/xray-node/docker-compose.yml \
    VERSION_PLACEHOLDER "$CLEAN_VER"

# env file
copy ./packaging/docker-compose/xray-node/example.env ./build/docker/xray-node/example.env

# config example folder
copy ./packaging/xray-node/config-example ./build/docker/xray-node/xray-node-config/example

# docker-compose archive
tar -czf docker-compose.xray-node.tar.gz -C ./build/docker/xray-node .

# ╔═══════════════════════════════════╗ #
# ║ nodeman with postgres and grafana ║ #
# ╚═══════════════════════════════════╝ #

# docker compose with replaced version placeholder
replace \
    ./packaging/docker-compose/xray-nodeman/docker-compose.yml \
    ./build/docker/xray-nodeman/docker-compose.yml \
    VERSION_PLACEHOLDER "$CLEAN_VER"

# env file
copy ./packaging/docker-compose/xray-nodeman/example.env ./build/docker/xray-nodeman/example.env


# copy prometheus config
copy ./packaging/docker-compose/xray-nodeman/datasource.yml ./build/docker/xray-nodeman/datasource.yml
copy ./packaging/docker-compose/xray-nodeman/prometheus.yml ./build/docker/xray-nodeman/prometheus.yml
       
# docker-compose archive
tar -czf docker-compose.xray-nodeman.tar.gz -C ./build/docker/xray-nodeman .

# ╔═══════════════════════════════════╗ #
# ║         nodeman standalone        ║ #
# ╚═══════════════════════════════════╝ #

# docker compose with replaced version placeholder
replace \
    ./packaging/docker-compose/xray-nodeman-standalone/docker-compose.yml \
    ./build/docker/xray-nodeman-standalone/docker-compose.yml \
    VERSION_PLACEHOLDER "$CLEAN_VER"
   
# docker-compose archive
tar -czf docker-compose.xray-nodeman-standalone.tar.gz -C ./build/docker/xray-nodeman-standalone .
