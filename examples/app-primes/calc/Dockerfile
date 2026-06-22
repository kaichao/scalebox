FROM python:3.14-slim

LABEL maintainer="kaichao"

COPY run.sh primes primes.py /app/bin/

COPY --from=scalebox.net/platform/agent /usr/local /usr/local
RUN mkdir -p /work
WORKDIR /work

ENTRYPOINT ["goagent"]
