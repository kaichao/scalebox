FROM debian:12-slim

LABEL maintainer="kaichao"

# install the newest version, rsync 3.2.7, zstd 1.5.5, openssh 9.3p1
RUN echo "deb http://deb.debian.org/debian testing main" > /etc/apt/sources.list.d/bookworm-testing.list \ 
    && apt-get update \
    && apt-get install -y rsync openssh-client zstd \
    && apt-get clean autoclean \
    && apt-get autoremove -y \
    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

COPY --from=hub.cstcloud.cn/scalebox/actuator /root/.ssh /root/.ssh

ENV SOURCE_CLUSTER= \
    TARGET_CLUSTER= \
    # 'local' / 'remote' / 'local#remote' / 'remote#local' 
    RELAY_OPTIONS= \
    # default value : 3
    ZSTD_CLEVEL= \
    KEEP_SOURCE_FILE=yes \
    # Break through the permissions in cifs-based file storage
    ENABLE_LOCAL_RELAY=

COPY *.sh /app/bin/

COPY --from=hub.cstcloud.cn/scalebox/agent /usr/local/ /usr/local/
RUN mkdir -p /work/.scalebox /app/bin && echo "PATH=/app/bin:\${PATH}" >> /root/.bashrc
WORKDIR /work
ENTRYPOINT [ "goagent" ]
