#!/bin/sh
set -eu

ARCH="${ARCH:-arm64}"

SERVICE="xray-node"

DUMMY_FILES="
build/$ARCH/xray-node/xray-node
build/$ARCH/xray-node/xray
build/$ARCH/xray-node/data/geoip.dat
build/$ARCH/xray-node/data/geosite.dat
packaging/xray-node/config-example/server.json
packaging/xray-node/config-example/client.json
"

OPERATIONS='
dpkg --install /build/package.deb
touch /etc/xray-node/config/real_client.json
touch /etc/xray-node/xray-node.env
touch /var/lib/xray-node/persistent/persistent-item
dpkg --remove xray-node
dpkg --install /build/package.deb
dpkg --purge xray-node
'


WATCH_DIRS="
/etc/$SERVICE
/usr/bin
/var/lib/$SERVICE
/lib/systemd/system
"

. ./packaging/deb-test/deb-test.sh

deb_test \
    "xray-node" \
    "packaging/xray-node/deb/nfpm.yaml" \
    "$DUMMY_FILES" \
    "$WATCH_DIRS" \
    "$OPERATIONS"