FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata

# Сборщик сам подставит нужную папку (amd64 или arm64) во время компиляции образа
ARG TARGETARCH
COPY ./build/${TARGETARCH}/xray-node /usr/bin/xray-node

# Копируем наш скрипт-обертку
COPY ./packaging/xray-node/docker/docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh

LABEL org.opencontainers.image.title="XRayMan Node"
LABEL org.opencontainers.image.description="View full help: docker compose run --rm xray-node --help"
LABEL org.opencontainers.image.vendor="XRayMan"

# Намертво фиксируем точку входа на наш скрипт
ENTRYPOINT ["/docker-entrypoint.sh"]