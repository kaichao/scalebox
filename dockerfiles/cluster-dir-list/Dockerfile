FROM debian:12-slim

LABEL maintainer="kaichao"

# install the newest version, rsync 3.2.7, openssh 9.3p1
RUN echo "deb http://deb.debian.org/debian testing main" > /etc/apt/sources.list.d/bookworm-testing.list \ 
    && apt-get update \
    && apt-get install -y rsync openssh-client zstd \
    && apt-get clean autoclean \
    && apt-get autoremove -y \
    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

COPY --from=hub.cstcloud.cn/scalebox/actuator /root/.ssh /root/.ssh

ENV SOURCE_CLUSTER= \
    REGEX_FILTER= \
    REGEX_2D_DATASET= \
    INDEX_2D_DATASET= \
    # 'yes'
    JUMP_SERVER_OPTION=

COPY run.sh /app/bin/
COPY list-files.sh /usr/local/bin/

COPY --from=hub.cstcloud.cn/scalebox/dir-list /usr/local/bin/get_2d_meta /usr/local/bin/

COPY --from=hub.cstcloud.cn/scalebox/agent /usr/local/ /usr/local/
RUN mkdir -p /work/.scalebox /app/bin && echo "PATH=/app/bin:\${PATH}" >> /root/.bashrc
WORKDIR /work
ENTRYPOINT [ "goagent" ]
