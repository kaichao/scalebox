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

FROM debian:12-slim

RUN \
    apt-get update \
    && apt-get install -y adminer \
    && apt-get clean autoclean \
    && apt-get autoremove -y \
    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

RUN echo "ServerName 127.0.0.1" >> /etc/apache2/apache2.conf \
    && a2enconf adminer \
    && sed -i 's/80/8080/' /etc/apache2/ports.conf

RUN \
    apt-get update \
    && apt-get install -y sqlite3 \
    && apt-get clean autoclean \
    && apt-get autoremove -y \
    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

COPY --from=hub.cstcloud.cn/scalebox/base /usr/local/sbin/ /usr/local/sbin/
RUN mkdir -p /work/.scalebox /app/bin && echo "PATH=/app/bin:\${PATH}" >> /root/.bashrc

COPY --from=0 /go/src/data-grouping /app/bin/data-grouping 
COPY adminer/entrypoint.sh /usr/local/bin/goagent.sh
COPY adminer/conf.php /etc/adminer/

ENV \
    ACTION_RUN=/app/bin/data-grouping \
    ALWAYS_RUNNING=yes \
    DATASET_PREFIX= \
    # 'string'/'integer', 'string' -> horizon-only
    COORD_TYPE=integer \
    # TRACE / DEBUG / INFO / WARN / ERROR / FATAL / PANIC
    LOG_LEVEL=INFO

# for adminer
EXPOSE 8080
WORKDIR /work

ENTRYPOINT [ "goagent.sh" ]
