ARG BUILDER_IMAGE
ARG BASE_IMAGE

FROM ${BUILDER_IMAGE}

FROM ${BASE_IMAGE}

LABEL maintainer="kaichao"

COPY --from=hub.cstcloud.cn/scalebox/agent /usr/local /usr/local

COPY --from=0 /src/ /app/bin/
COPY run.sh /app/bin/

ENV LENGTH=10000

RUN mkdir -p /work
WORKDIR /work

ENTRYPOINT ["goagent"]
