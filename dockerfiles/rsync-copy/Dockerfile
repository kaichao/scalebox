FROM golang:1.20.4

COPY jump_server/jump_servers.go /src/
RUN cd /src && go build jump_servers.go && strip jump_servers

FROM hub.cstcloud.cn/scalebox/agent

LABEL maintainer="kaichao"

# install the newest version, rsync 3.2.7, zstd 1.5.4
RUN echo "deb http://deb.debian.org/debian testing main" > /etc/apt/sources.list.d/bullseye-testing.list \ 
    && apt-get update \
    && apt-get install -y rsync openssh-client zstd \
    && apt-get clean autoclean \
    && apt-get autoremove -y \
    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

COPY --from=hub.cstcloud.cn/scalebox/actuator /root/.ssh /root/.ssh
COPY --from=hub.cstcloud.cn/scalebox/list-dir /app/bin/url_parser /app/bin/ 
COPY --from=0 /src/jump_servers /app/bin/ 

WORKDIR /work

ENV SOURCE_URL= \
    TARGET_URL= \
    RSYNC_PASSWORD= \
    # user1@host1,user2@host2
    JUMP_SERVERS= \
    ZSTD_CLEVEL=3 \
    ENABLE_ZSTD= \
    # Break through the permissions in cifs-based file storage
    ENABLE_LOCAL_RELAY= \
    SERVER_CONTAINERIZED=

COPY *.sh /app/bin/
