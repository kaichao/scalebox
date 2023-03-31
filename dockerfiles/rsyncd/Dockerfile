FROM hub.cstcloud.cn/scalebox/rsync-agent

LABEL maintainer="kaichao"

RUN apt-get update \
    && apt-get install -y openssh-server sudo \
    && apt-get clean autoclean \
    && apt-get autoremove -y \
    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

RUN useradd -rm -d /home/ubuntu -s /bin/bash -g root -G sudo -u 1000 test \
    && usermod -aG sudo test \
    && echo 'test:test' | chpasswd

RUN service ssh start

ENV TRANSPORT_TYPE=rsync \
    RSYNC_PORT=873 \
    SSH_PORT=2222

ENV ACTION_RUN=/usr/local/bin/run.sh

COPY rootfs/ /

COPY --from=hub.cstcloud.cn/scalebox/rsync-base /root/.ssh/id_rsa.pub /root/.ssh/authorized_keys

# CMD run.sh
