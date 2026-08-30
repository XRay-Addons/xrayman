FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata

# Docker Buildx сам подставит нужную папку во время multi-arch сборки
ARG TARGETARCH
COPY ./downloads/xrayman-nodeman-linux-${TARGETARCH}.tar.gz /tmp/
RUN tar -xzf /tmp/xrayman-nodeman-linux-${TARGETARCH}.tar.gz -C /usr/bin/ && rm /tmp/*.tar.gz

ENTRYPOINT ["/usr/bin/xray-nodeman"]