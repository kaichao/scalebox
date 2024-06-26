FROM golang:1.20.6

ARG GOPROXY=https://goproxy.cn

WORKDIR /go/src

COPY go.* /go/src/
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY cmd /go/src/cmd

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build ./cmd/url_parser && strip url_parser \
    && go build ./cmd/get_2d_meta && strip get_2d_meta \
    && go build ./cmd/jump_servers && strip jump_servers

FROM debian:11-slim

LABEL maintainer="kaichao"

# the newest curlftpfs may not be compatible, the following is OK.
# $ curlftpfs --version
# curlftpfs 0.9.2 libcurl/7.74.0 fuse/2.9
RUN \
    apt-get update \
    && apt-get install -y curlftpfs \
    && apt-get clean autoclean \
    && apt-get autoremove -y \
    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# install the newest version, rsync 3.2.6, zstd 1.5.2
RUN echo "deb http://deb.debian.org/debian testing main" > /etc/apt/sources.list.d/bullseye-testing.list \
    && apt-get update \
    && apt-get install -y rsync openssh-client \
    && apt-get clean autoclean \
    && apt-get autoremove -y \
    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

COPY --from=hub.cstcloud.cn/scalebox/actuator /root/.ssh /root/.ssh
COPY --from=0 /go/src/url_parser /go/src/get_2d_meta /go/src/jump_servers /usr/local/bin/ 

RUN mkdir -p /remote

ENV SOURCE_URL= \
    REGEX_FILTER= \
    REGEX_2D_DATASET= \
    INDEX_2D_DATASET= \
    RSYNC_PASSWORD=

COPY run.sh /app/share/bin/
COPY list-files.sh /usr/local/bin/

COPY --from=hub.cstcloud.cn/scalebox/agent /usr/local /usr/local
RUN mkdir -p /work/.scalebox /app/bin && echo "PATH=/app/bin:\${PATH}" >> /root/.bashrc
WORKDIR /work
ENTRYPOINT [ "goagent" ]
