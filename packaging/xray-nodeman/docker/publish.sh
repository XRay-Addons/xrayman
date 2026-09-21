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

# docker compose with replaced version placeholder
replace \
    ./packaging/xray-nodeman/docker/docker-compose.yml \
    ./build/docker/xray-nodeman/docker-compose.yml \
    VERSION_PLACEHOLDER "$CLEAN_VER"


# docker-compose archive
tar -czf docker-compose.xray-nodeman.tar.gz -C ./build/docker/xray-nodeman .
