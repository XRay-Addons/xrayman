FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata

# Копируем уже скомпилированный бинарник текущей архитектуры матрицы
COPY ./build/xray-nodeman/xray-nodeman /usr/bin/xray-nodeman

ENTRYPOINT ["/usr/bin/xray-nodeman"]
