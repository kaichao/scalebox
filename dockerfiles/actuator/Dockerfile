FROM hub.cstcloud.cn/scalebox/actuator

RUN \
    echo "auto keygen in actuator ..." ; \
    rm -f /root/.ssh/id_rsa ; \
    ssh-keygen -q -N "" -f /root/.ssh/id_rsa ;
