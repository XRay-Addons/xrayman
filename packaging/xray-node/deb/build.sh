#!/bin/bash
set -e
 
VERSION="${1:-dev}"
BUILD_ROOT="${2:-./build}"
ARCH="${3:-amd64}"                 # без ведущей точки, в отличие от старого скрипта
FORMAT="${4:-deb}"                 # deb | rpm | apk — что угодно, что умеет nfpm
 
export VERSION
export ARCH
export BUILD_DIR="$BUILD_ROOT/$ARCH/xray-node"
 
# Сохраняем старое имя и расположение файла, чтобы ничего не пришлось
# менять в publish/деплой-скриптах: раньше dpkg-deb клал .deb рядом
# с директорией "xray-node-${ARCH}" в текущем рабочем каталоге, т.е.
# итоговый путь был "./xray-node-${ARCH}.deb".
OUT_FILE="./xray-node-${ARCH}.${FORMAT}"
 
nfpm package \
    --config ./packaging/xray-node/deb/nfpm.yaml \
    --target "$OUT_FILE" \
    --packager "$FORMAT"
 
echo "Built: $OUT_FILE"