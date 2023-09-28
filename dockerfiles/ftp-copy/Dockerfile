FROM debian:11-slim

LABEL maintainer="kaichao"

RUN echo "deb http://deb.debian.org/debian testing main" > /etc/apt/sources.list.d/bullseye-testing.list \
    && apt-get update \
    && apt-get install -y curlftpfs lftp \
    && apt-get clean autoclean \
    && apt-get autoremove -y \
    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

COPY --from=hub.cstcloud.cn/scalebox/list-dir /usr/local/bin/url_parser /usr/local/bin/ 

RUN mkdir -p /local /work

COPY *.sh /app/bin/

ENV SOURCE_URL= \
    TARGET_URL= \
    ACTION= \
    ENABLE_RECHECK_SIZE=yes \
    NUM_PGET_CONN=4 \
    ENABLE_LOCAL_RELAY=no \
    RAM_DISK_GB= 

COPY --from=hub.cstcloud.cn/scalebox/base /usr/local/sbin/ /usr/local/sbin/
RUN mkdir -p /work /app/bin && echo "PATH=/app/bin:\${PATH}" >> /root/.bashrc
WORKDIR /work
ENTRYPOINT [ "goagent" ]
