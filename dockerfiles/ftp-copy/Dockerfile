FROM hub.cstcloud.cn/scalebox/agent

LABEL maintainer="kaichao"

RUN echo "deb http://deb.debian.org/debian testing main" > /etc/apt/sources.list.d/bullseye-testing.list \
    && apt-get update \
    && apt-get install -y curlftpfs lftp \
    && apt-get clean autoclean \
    && apt-get autoremove -y \
    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

COPY --from=hub.cstcloud.cn/scalebox/list-dir /app/bin/url_parser /app/bin/ 

RUN mkdir -p /local /work

COPY *.sh /app/bin/

ENV SOURCE_URL= \
    TARGET_URL= \
    ACTION= \
    ENABLE_RECHECK_SIZE=yes \
    NUM_PGET_CONN=4 \
    ENABLE_LOCAL_RELAY=no \
    RAM_DISK_GB= 
