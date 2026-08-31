FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata

# Сборщик сам подставит нужную папку (amd64 или arm64) во время компиляции образа
ARG TARGETARCH
COPY ./dist/${TARGETARCH}/xray-nodeman /usr/bin/xray-nodeman

# Копируем наш скрипт-обертку
COPY ./packaging/docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh

LABEL org.opencontainers.image.title="XRayMan Node Manager"
LABEL org.opencontainers.image.description="View full help: docker compose run --rm xray-nodeman --help"
LABEL org.opencontainers.image.vendor="XRayMan"

# Намертво фиксируем точку входа на наш скрипт
ENTRYPOINT ["/docker-entrypoint.sh"]