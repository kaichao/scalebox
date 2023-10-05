FROM golang:1.20.6

LABEL maintainer="kaichao"

ARG GOPROXY=https://goproxy.io

WORKDIR /go/src
COPY go.* /go/src/
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY *.go /go/src/
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    # glibc included.
    # CGO_ENABLED=0 \
    go build 

FROM hub.cstcloud.cn/scalebox/agent

COPY --from=0 /go/src/data-grouping /app/bin/data-grouping 

RUN mkdir -p /work/.data-grouping
ENV \
    ACTION_RUN=/app/bin/data-grouping \
    MESSAGE_FILE=/work/messages.txt \
    DATASET_FILE=/work/.data-grouping/datasets.txt \
    SQLITE_FILE=/work/.data-grouping/my.db \
    ALWAYS_RUNNING=yes \
    # TRACE / DEBUG / INFO / WARN / ERROR / FATAL / PANIC
    LOG_LEVEL=INFO
