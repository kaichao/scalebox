FROM python:3.11-slim

LABEL maintainer="kaichao"

COPY --from=hub.cstcloud.cn/scalebox/base /usr/local/sbin /usr/local/sbin
COPY run.sh primes primes.py /app/bin/

ENV LENGTH=10000

RUN mkdir -p /work
WORKDIR /work

ENTRYPOINT ["goagent"]
