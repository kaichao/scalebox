clusters ?= local bio-down

include defs.mk

# Get Local IPv4 addr
# ip route get 1.2.3.4 | awk '{print $7}'
# hostname -I | awk '{print $1}'
# hostname -i
# local_ip_index=4
# local_addr=$$(hostname -I| cut -d ' ' -f ${local_ip_index})
ifeq ($(origin LOCAL_ADDR), undefined)
	ifeq ($(origin LOCAL_IP_INDEX), undefined)
	    LOCAL_ADDR := $$(hostname -i)
	else
	    LOCAL_ADDR := $$(hostname -I| cut -d ' ' -f ${LOCAL_IP_INDEX})
    endif
endif

default:

all: reset import

reset: down
	@# docker ps -a|grep -v CONTAINER|awk '{print $$1}'| while read line; do docker rm -f $$line;done
	@ rm -rf ${HOME}/.scalebox/log/cli
	@ mkdir -p ${HOME}/.scalebox/log/cli
	@ docker volume rm -f cluster_pgdata
	@ docker-compose up -d
	@ sleep 25

import:
	@ echo "LOCAL_ADDR is "$(LOCAL_ADDR)
	@ for c in $(clusters); do \
		LOCAL_ADDR=$(LOCAL_ADDR) scalebox app create $$c/mycluster.yaml; \
	done

up:
	@ docker-compose up -d

down:
	@ docker-compose down

clean:
	@ pdsh -w $(NODES) 'docker rm -f $$(docker ps -qa)' | dshbak -c
	@ echo "================"
	@ docker rm -f $$(docker ps -qa)

list:
	@ pdsh -w $(NODES) 'docker  ps -a' | dshbak -c
	@ echo "================"
	@ docker ps -a

actuator-passwordless:
	@ mkdir -p $(SHARED_DIR) && cp id_rsa.pub $(SHARED_DIR)
	@ pdsh -w $(NODES) 'cat $(SHARED_DIR)/id_rsa.pub >> ${HOME}/.ssh/authorized_keys' | dshbak -c
	@ rm -f $(SHARED_DIR)/id_rsa.pub
