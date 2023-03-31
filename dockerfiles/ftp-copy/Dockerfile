FROM kaichao/ftp-copy

LABEL maintainer="kaichao"

ENV REMOTE_URL= \
    LOCAL_ROOT=/ \
    NUM_PGET_CONN=4 \
    ACTION=PUSH \
    ENABLE_RECHECK_SIZE=yes \
    ENABLE_LOCAL_RELAY= \
    RAM_DISK_GB=

COPY --from=hub.cstcloud.cn/scalebox/base /usr/local/sbin /usr/local/sbin

ENTRYPOINT ["goagent"]
