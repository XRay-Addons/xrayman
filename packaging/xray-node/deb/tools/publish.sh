#!/bin/bash
. ./packaging/common/utils.sh

set -e

CLEAN_VER="${1:-0.0.0}"

# docker compose with replaced version placeholder
replace \
    ./packaging/xray-node/docker/docker-compose.yml \
    ./build/docker/xray-node/docker-compose.yml \
    VERSION_PLACEHOLDER "$CLEAN_VER"

# config example folder
copy ./packaging/xray-node/config ./build/docker/xray-node/xray-node-config/example

# docker-compose archive
tar -czf docker-compose.xray-node.tar.gz -C ./build/docker/xray-node .

          # 1. Создаем целевые папки (без этого sed упадет с ошибкой)
          # mkdir -p ./dist/docker
          # mkdir -p ./dist/docker/xray-nodeman
          # mkdir -p ./dist/docker/xray-nodeman/all-in-one
          # mkdir -p ./dist/docker/xray-nodeman/standalone
          # mkdir -p ./dist/docker/xray-node

          # 2. Заменяем плейсхолдер версии на реальный тег релиза
          # sed "s/VERSION_PLACEHOLDER/$CLEAN_VER/g" ./packaging/xray-nodeman/docker/all-in-one/docker-compose.yml > ./dist/docker/xray-nodeman/all-in-one/docker-compose.yml
          # sed "s/VERSION_PLACEHOLDER/$CLEAN_VER/g" ./packaging/xray-nodeman/docker/standalone/docker-compose.yml > ./dist/docker/xray-nodeman/standalone/docker-compose.yml
          # sed "s/VERSION_PLACEHOLDER/$CLEAN_VER/g" ./packaging/xray-node/docker/docker-compose.yml > ./dist/docker/xray-node/docker-compose.yml

          # 3. (Опционально) Скопируйте сюда дополнительные файлы конфигурации, если они есть
          # cp ./packaging/xray-nodeman/docker/all-in-one/datasource.yml ./dist/docker/xray-nodeman/all-in-one/datasource.yml
          # cp ./packaging/xray-nodeman/docker/all-in-one/prometheus.yml ./dist/docker/xray-nodeman/all-in-one/prometheus.yml
          # cp -r ./packaging/xray-node/config ./dist/docker/xray-node/xray-node-config

          # 4. Архиривуем папки (флаг -C временно меняет директорию, чтобы внутри архива не было лишней вложенности)
          # tar -czf docker-compose.xray-nodeman.all-in-one.tar.gz -C ./dist/docker/xray-nodeman/all-in-one .
          # tar -czf docker-compose.xray-nodeman.standalone.tar.gz -C ./dist/docker/xray-nodeman/standalone .
          # tar -czf docker-compose.xray-node.tar.gz -C ./dist/docker/xray-node .

