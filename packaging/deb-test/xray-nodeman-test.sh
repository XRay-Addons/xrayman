#!/bin/sh
set -eu

ARCH="${ARCH:-arm64}"

SERVICE="xray-nodeman"

DUMMY_FILES="
build/$ARCH/xray-nodeman/xray-nodeman
"

OPERATIONS='
dpkg --install /build/package.deb
touch /etc/xray-nodeman/xray-nodeman.env
dpkg --remove xray-nodeman
dpkg --install /build/package.deb
dpkg --purge xray-nodeman
'


WATCH_DIRS="
/etc/$SERVICE
/usr/bin
/var/lib/$SERVICE
/lib/systemd/system
"

. ./packaging/deb-test/deb-test.sh

deb_test \
    "xray-nodeman" \
    "packaging/xray-nodeman/deb/nfpm.yaml" \
    "$DUMMY_FILES" \
    "$WATCH_DIRS" \
    "$OPERATIONS"