FROM golang:1.20.3

ARG GOPROXY=https://goproxy.cn

COPY . /src/
RUN --mount=type=cache,target=/go/pkg \
    cd /src && go build main.go && strip main

FROM hub.cstcloud.cn/scalebox/agent

COPY --from=0 /src/main /app/bin/cron

ENV ACTION_RUN=/app/bin/cron

ENTRYPOINT ["goagent"]
