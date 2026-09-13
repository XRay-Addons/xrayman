#!/bin/bash
set -e

VERSION="${1:-dev}"
BUILD_ROOT="${2:-./build}"

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

    sed "${args[@]}" "$src" > "$dst"
}

for ARCH in amd64 arm64; do
    DIR="xray-node-${ARCH}"
    BUILD_DIR="$BUILD_ROOT/${ARCH}/xray-node"
    PKG_DIR="./packaging/xray-node"

    rm -rf "$DIR"

    # copy build results
    copy "$BUILD_DIR/xray-node" "$DIR/usr/bin/xray-node"
    copy "$BUILD_DIR/data" "$DIR/var/lib/xray-node/data"
    copy "$BUILD_DIR/xray" "$DIR/usr/bin/xray"

    # copy packaging stuff
    copy "$PKG_DIR/config" "$DIR/etc/xray-node/config"
    copy "$PKG_DIR/deb/xray-node.service" "$DIR/lib/systemd/system/xray-node.service"
    copy "$PKG_DIR/deb/xray-node.example.env" "$DIR/etc/xray-node/xray-node.example.env"
    
    # copy deb scripts
    copy "$PKG_DIR/deb/postinstall.sh" "$DIR/DEBIAN/postinst" +x
    copy "$PKG_DIR/deb/prerm.sh" "$DIR/DEBIAN/prerm" +x
    copy "$PKG_DIR/deb/postrm.sh" "$DIR/DEBIAN/postrm" +x

    # deb template - set version and arch
    replace \
        "$PKG_DIR/deb/control.template" \
        "$DIR/DEBIAN/control" \
        VERSION_PLACEHOLDER "$VERSION" \
        ARCH_PLACEHOLDER "$ARCH"

    dpkg-deb --build "$DIR"
done