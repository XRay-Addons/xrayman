FROM alpine:3.22

RUN apk add --no-cache ca-certificates

ARG TARGETARCH

COPY xray-nodeman-${TARGETARCH} /usr/local/bin/xray-nodeman

ENTRYPOINT ["/usr/local/bin/xray-nodeman"]