FROM golang:1.20.4

COPY url_parser.go /src/
RUN cd /src && go build url_parser.go && strip url_parser

FROM hub.cstcloud.cn/scalebox/agent

LABEL maintainer="kaichao"

# the newest curlftpfs may not compatible, the following is OK.
# $ curlftpfs --version
# curlftpfs 0.9.2 libcurl/7.74.0 fuse/2.9
RUN \
    apt-get update \
    && apt-get install -y curlftpfs \
    && apt-get clean autoclean \
    && apt-get autoremove -y \
    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# install the newest version, rsync 3.2.6, zstd 1.5.2
RUN echo "deb http://deb.debian.org/debian testing main" > /etc/apt/sources.list.d/bullseye-testing.list \
    && apt-get update \
    && apt-get install -y rsync openssh-client \
    && apt-get clean autoclean \
    && apt-get autoremove -y \
    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

COPY --from=hub.cstcloud.cn/scalebox/actuator /root/.ssh /root/.ssh
COPY --from=0 /src/url_parser /app/bin/ 

RUN mkdir -p /remote

ENV SOURCE_URL= \
    REGEX_FILTER= \
    RSYNC_PASSWORD=

COPY *.sh /app/bin/
