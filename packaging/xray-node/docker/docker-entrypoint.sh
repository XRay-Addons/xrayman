#!/bin/sh
set -e

# Генерация уникальных ключей при первом запуске - для примеров
XRAY_CONFIG_DIR=/xray-node-config/
if [ -f "$XRAY_CONFIG_DIR/xray_keygen.tmp.sh" ]; then
    echo "Generate config example key."

    "$XRAY_CONFIG_DIR/xray_keygen.tmp.sh" /usr/bin/xray-node/xray \
        "$XRAY_CONFIG_DIR/xray_server.example.json:$XRAY_CONFIG_DIR/xray_client.example.json" \
        X25519_PRIVATE_KEY X25519_PUBLIC_KEY

    rm "$XRAY_CONFIG_DIR/xray_keygen.tmp.sh"
fi

# Рисуем красивую рамку-подсказку при каждом старте контейнера
echo "╔═══════════════════════════════════════════════════════════════════════╗"
echo "║               XRAYMAN DOCKER CONTAINER STARTING...                    ║"
echo "╠═══════════════════════════════════════════════════════════════════════╣"
echo "║                                                                       ║"
echo "║  * view actual help and env variables:                                ║"
echo "║    docker compose run --rm xray-node --help                           ║"
echo "║                                                                       ║"
echo "║  * create xray server and client configs in xray-node-config          ║"
echo "║    (see example inside)                                               ║"
echo "║                                                                       ║"
echo "║  * start xray-node:                                                   ║"
echo "║    docker compose up -d                                               ║"
echo "║                                                                       ║"
echo "║  * restart xray-node:                                                 ║"
echo "║    docker compose down && docker compose up -d                        ║"
echo "║                                                                       ║"
echo "║  * view logs:                                                         ║"
echo "║    docker compose logs -n 50 -f xray-node                             ║"
echo "║                                                                       ║"
echo "║  * node access key will be stored in xray-node-persistent:            ║"
echo "║    cat xray-node-persistent/access.json.                              ║"
echo "║                                                                       ║"
echo "╚═══════════════════════════════════════════════════════════════════════╝"
echo ""

# Магия: передаем управление реальному приложению, прокидывая все аргументы ($@)
exec /usr/bin/xray-node/xray-node "$@"