FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata

# Сборщик сам подставит нужную папку (amd64 или arm64) во время компиляции образа
ARG TARGETARCH
COPY ./dist/${TARGETARCH}/xray-nodeman /usr/bin/xray-nodeman

ENTRYPOINT ["/usr/bin/xray-nodeman"]
