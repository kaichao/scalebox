FROM hub.cstcloud.cn/scalebox/rsync-agent

LABEL maintainer="kaichao"

ENV SOURCE_URL= \
    TARGET_URL= \
    RSYNC_PASSWORD= \
    ZSTD_CLEVEL=3 \
    ENABLE_ZSTD= \
    # Break through the permissions of some file storage
    ENABLE_LOCAL_RELAY= \
    SERVER_CONTAINERIZED=


COPY *.sh /app/bin/
