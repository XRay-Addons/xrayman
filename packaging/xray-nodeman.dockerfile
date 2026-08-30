FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata
COPY ./build/xray-nodeman/xray-nodeman /usr/bin/xrayman-nodeman
ENTRYPOINT ["/usr/bin/xrayman-nodeman"]