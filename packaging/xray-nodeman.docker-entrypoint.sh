#!/bin/sh
set -e

# Рисуем красивую рамку-подсказку при каждом старте контейнера
echo "╔═══════════════════════════════════════════════════════════════════════╗"
echo "║               XRAYMAN DOCKER CONTAINER STARTING...                    ║"
echo "╠═══════════════════════════════════════════════════════════════════════╣"
echo "║                                                                       ║"
echo "║  * view actual help and env variables:                                ║"
echo "║    docker compose run --rm xray-nodeman --help                        ║"
echo "║                                                                       ║"
echo "║  * create and edit .env file and set up env variables:                ║"
echo "║    sudo nano .env                                                     ║"
echo "║                                                                       ║"
echo "║  * start xray-nodeman:                                                ║"
echo "║    docker compose up -d                                               ║"
echo "║                                                                       ║"
echo "║  * view logs:                                                         ║"
echo "║    docker compose logs -n 50 -f xray-nodeman                          ║"
echo "║                                                                       ║"
echo "╚═══════════════════════════════════════════════════════════════════════╝"
echo ""

# Магия: передаем управление реальному приложению, прокидывая все аргументы ($@)
exec /usr/bin/xray-nodeman "$@"