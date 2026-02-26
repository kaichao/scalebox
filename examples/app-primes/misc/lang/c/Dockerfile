FROM gcc:15-trixie

COPY . /src/

WORKDIR /src/

# static link
RUN gcc -std=c99 -D CONTAINER_APP -Iinclude -static -fopenmp -O3 primes.c -o primes

# FROM debian:13-slim

# LABEL maintainer="kaichao"

# # RUN apt-get update \
# #     && apt-get install -y libgomp1 \
# #     && apt-get clean \
# #     && rm -rf /var/lib/apt/lists/*
# COPY --from=hub.cstcloud.cn/scalebox/agent:latest /usr/local/sbin /usr/local/sbin

# ENV OMP_NUM_THREADS= \
#     GROUP_SIZE=10000
# COPY --from=0 /src/primes /usr/local/bin
# COPY run.sh /app/bin/

# ENTRYPOINT ["goagent"]
