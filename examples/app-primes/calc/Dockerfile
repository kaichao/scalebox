FROM python:3.11-slim

LABEL maintainer="kaichao"

COPY run.sh primes primes.py /app/bin/

ENV LENGTH=10000

COPY --from=hub.cstcloud.cn/scalebox/agent /usr/local /usr/local
RUN mkdir -p /work
WORKDIR /work

ENTRYPOINT ["goagent"]
