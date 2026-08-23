FROM scalebox.net/platform/agent

COPY run.sh check.sh /app/bin/

ENV HEAD_MODULE=vtask-head \
    VTASK_MODE=HOST-BOUND

ENTRYPOINT ["goagent"]
