FROM debian:13-slim

LABEL maintainer="kaichao"

# install the newest version, 
#   rsync 3.4.1
#   zstd 1.5.7
#   OpenSSH_10.0p2 Debian-7, OpenSSL 3.5.1 1 Jul 2025
#   pv 1.9.31
# RUN echo "deb http://deb.debian.org/debian testing main" > /etc/apt/sources.list.d/bookworm-testing.list \ 
RUN \
    apt-get update \
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
    SOURCE_URL= \
    TARGET_URL= \
    # LOCAL / SSH / RSYNC_OVER_SSH /RSYNC
    SOURCE_MODE= \
    TARGET_MODE= \
    SOURCE_JUMP= \
    TARGET_JUMP= \
    TARGET_FILE_MODE=0644 \
    ENABLE_TARGET_TOUCH=no \
    ZSTD_CLEVEL= \
    KEEP_SOURCE_FILE=yes \
    # MD-5 / SHA-1 / SHA-256 / SHA-512 / RIPEMD-160
    CHECKSUM_ALGO= \
    RSYNC_BWLIMIT= \
    RSYNC_PASSWORD= \
    RSYNCD_MODULE=/root \
    WORK_DIR=

COPY code/*.sh /app/share/bin/

COPY --from=hub.cstcloud.cn/scalebox/agent /usr/local/ /usr/local/
RUN mkdir -p /work/.scalebox /app/bin && echo "PATH=/app/bin:\${PATH}" >> /root/.bashrc
WORKDIR /work
ENTRYPOINT [ "entrypoint.sh" ]
