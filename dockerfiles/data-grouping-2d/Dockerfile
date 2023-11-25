FROM golang:1.20.6

LABEL maintainer="kaichao"

ARG GOPROXY=https://goproxy.io

WORKDIR /go/src
COPY go.* /go/src/
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY *.go /go/src/
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    # glibc included.
    # CGO_ENABLED=0 \
    go build 

FROM hub.cstcloud.cn/scalebox/agent

RUN \
    apt-get update \
    && apt-get install -y adminer \
    && apt-get clean autoclean \
    && apt-get autoremove -y \
    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

RUN echo "ServerName 127.0.0.1" >> /etc/apache2/apache2.conf \
    && a2enconf adminer

COPY --from=0 /go/src/data-grouping /app/bin/data-grouping 
COPY adminer/entrypoint.sh /usr/local/bin/goagent.sh
COPY adminer/conf.php /etc/adminer/

RUN mkdir -p /work/.data-grouping
ENV \
    ACTION_RUN=/app/bin/data-grouping \
    MESSAGE_FILE=/work/messages.txt \
    DATASET_FILE=/work/.data-grouping/datasets.txt \
    SQLITE_FILE=/work/.data-grouping/my.db \
    ALWAYS_RUNNING=yes \
    DATASET_PREFIX= \
    # 'string'/'integer', 'string' -> horizon-only
    COORD_TYPE= \
    # TRACE / DEBUG / INFO / WARN / ERROR / FATAL / PANIC
    LOG_LEVEL=INFO

# for adminer
EXPOSE 80

ENTRYPOINT [ "goagent.sh" ]
