FROM debian:12-slim

LABEL maintainer="kaichao"

# install the newest version, rsync 3.2.7, zstd 1.5.5, openssh 9.3p1
RUN echo "deb http://deb.debian.org/debian testing main" > /etc/apt/sources.list.d/bookworm-testing.list \ 
    && apt-get update \
    && apt-get install -y rsync openssh-client zstd pv \
    && apt-get clean autoclean \
    && apt-get autoremove -y \
    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# jq for json parsing in bash
RUN rm -f /etc/apt/apt.conf.d/docker-clean; \
    echo 'Binary::apt::APT::Keep-Downloaded-Packages "true";' > /etc/apt/apt.conf.d/keep-cache
RUN --mount=type=cache,target=/var/cache/apt  \
    --mount=type=cache,target=/var/lib/apt \
    apt update \
    && apt-get install -y jq \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

COPY --from=hub.cstcloud.cn/scalebox/actuator /root/.ssh /root/.ssh

ENV \
    WITH_HEADERS=yes \
    SOURCE_URL= \
    TARGET_URL= \
    SOURCE_JUMP_SERVERS= \
    TARGET_JUMP_SERVERS= \
    RSYNC_OVER_SSH= \
    ZSTD_CLEVEL= \
    KEEP_SOURCE_FILE=yes \
    WORK_DIR=

COPY remote-copy/run.sh /app/share/bin/
COPY functions.sh /app/share/bin/

COPY --from=hub.cstcloud.cn/scalebox/agent /usr/local/ /usr/local/
RUN mkdir -p /work/.scalebox /app/bin && echo "PATH=/app/bin:\${PATH}" >> /root/.bashrc
WORKDIR /work
ENTRYPOINT [ "entrypoint.sh" ]
