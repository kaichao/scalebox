FROM debian:11-slim

LABEL maintainer="kaichao"

# install the newest version, rsync 3.2.7, zstd 1.5.4
RUN echo "deb http://deb.debian.org/debian testing main" > /etc/apt/sources.list.d/bullseye-testing.list \ 
    && apt-get update \
    && apt-get install -y rsync openssh-client zstd \
    && apt-get clean autoclean \
    && apt-get autoremove -y \
    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

RUN apt-get update \
    && apt-get install -y openssh-server sudo \
    && apt-get clean autoclean \
    && apt-get autoremove -y \
    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

RUN useradd -rm -d /home/ubuntu -s /bin/bash -g root -G sudo -u 1000 test \
    && usermod -aG sudo test \
    && echo 'test:test' | chpasswd

RUN service ssh start

ENV MODE=rsyncd \
    RSYNC_PORT=873 \
    SSH_PORT=2222

ENV ACTION_RUN=/usr/local/bin/run.sh

COPY rootfs/ /

COPY --from=hub.cstcloud.cn/scalebox/actuator /root/.ssh/id_rsa.pub /root/.ssh/authorized_keys
COPY --from=hub.cstcloud.cn/scalebox/list-dir /app/bin/url_parser /app/bin/ 

COPY --from=hub.cstcloud.cn/scalebox/base /usr/local/sbin/ /usr/local/sbin/
RUN mkdir -p /work /app/bin && echo "PATH=/app/bin:\${PATH}" >> /root/.bashrc
WORKDIR /work
ENTRYPOINT [ "goagent" ]
