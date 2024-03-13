FROM debian:11-slim

LABEL maintainer="kaichao"

RUN echo "deb http://deb.debian.org/debian testing main" > /etc/apt/sources.list.d/bullseye-testing.list \
    && apt-get update \
    && apt-get install -y curlftpfs lftp \
    && apt-get clean autoclean \
    && apt-get autoremove -y \
    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

RUN mkdir -p /local /work

COPY *.sh /app/bin/

ENV SOURCE_URL= \
    TARGET_URL= \
    ENABLE_RECHECK_SIZE=yes \
    ENABLE_RECHECK_PUSH= \
    NUM_PGET_CONN=4 \
    ENABLE_LOCAL_RELAY=no

COPY --from=hub.cstcloud.cn/scalebox/agent /usr/local/ /usr/local/
RUN mkdir -p /work/.scalebox /app/bin && echo "PATH=/app/bin:\${PATH}" >> /root/.bashrc
WORKDIR /work
ENTRYPOINT [ "goagent" ]
