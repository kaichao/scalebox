FROM golang:1.19.2

COPY . /src/
RUN cd /src && go build url_parser.go && strip url_parser

FROM debian:11-slim AS agent-base

LABEL maintainer="kaichao"

# install the newest version, rsync 3.2.6, zstd 1.5.2
RUN echo "deb http://deb.debian.org/debian testing main" > /etc/apt/sources.list.d/bullseye-testing.list \ 
    && apt-get update \
    && apt-get install -y rsync openssh-client zstd \
    && apt-get clean autoclean \
    && apt-get autoremove -y \
    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

COPY --from=hub.cstcloud.cn/scalebox/base /usr/local/sbin /usr/local/sbin
COPY --from=hub.cstcloud.cn/scalebox/actuator /root/.ssh /root/.ssh
COPY --from=0 /src/url_parser /app/bin/ 

RUN mkdir -p /work /app/bin \
    && echo "PATH=/app/bin:\${PATH}" >> /root/.bashrc
WORKDIR /work

FROM agent-base

ENV \
    REMOTE_HOST=localhost \
    REMOTE_PORT=22 \
    REMOTE_USER=root \
    # for containerized remote rsyncd
    REMOTE_ROOT=/local \
    LOCAL_ROOT=/local \
    RSYNC_PASSWORD=

ENTRYPOINT ["goagent"]
