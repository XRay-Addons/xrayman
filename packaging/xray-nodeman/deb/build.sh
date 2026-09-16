#!/bin/bash

. ./packaging/common/utils.sh

set -e

VERSION="${1:-dev}"
BUILD_ROOT="${2:-./build}"
ARCH="${3:-.amd64}"

DIR="xray-nodeman-${ARCH}"
BUILD_DIR="$BUILD_ROOT/${ARCH}/xray-nodeman"
PKG_DIR="./packaging/xray-nodeman"

rm -rf "$DIR"

# copy build results
copy "$BUILD_DIR/xray-nodeman" "$DIR/usr/bin/xray-nodeman"

# copy packaging stuff
copy "$PKG_DIR/deb/xray-nodeman.service" "$DIR/lib/systemd/system/xray-nodeman.service"
copy "$PKG_DIR/deb/xray-nodeman.example.env" "$DIR/etc/xray-node/xray-nodeman.example.env"

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

dpkg-deb --root-owner-group --build "$DIR"