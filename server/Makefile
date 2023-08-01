clusters ?= local inline-cluster external-cluster

include defs.mk

# Get Local IPv4 addr
# ip route get 1.2.3.4 | awk '{print $7}'
# hostname -I | awk '{print $1}'
# hostname -i
# local_ip_index=4
# local_addr=$$(hostname -I| cut -d ' ' -f ${local_ip_index})

UNAME := $(shell uname)

ifeq ($(origin LOCAL_ADDR), undefined)
	ifeq ($(UNAME), Linux)
		ifeq ($(origin LOCAL_IP_INDEX), undefined)
		    LOCAL_ADDR := $$(hostname -i)
		else
	    	LOCAL_ADDR := $$(hostname -I| cut -d ' ' -f ${LOCAL_IP_INDEX})
    	endif
	else ifeq ($(UNAME), Darwin)
		LOCAL_ADDR := $$(ifconfig|grep "inet "|grep -v 127.0.0.1|head -1|cut -d ' ' -f 2)
	else
		$(info This is neither a Linux nor a macOS system)
	endif

endif

default:
	@ echo "LOCAL_ADDR is "$(LOCAL_ADDR)

all: reset import

reset: down
	@ rm -rf ${HOME}/.scalebox/log/cli
	@ mkdir -p ${HOME}/.scalebox/log/cli
	@ docker volume rm -f server_pgdata
	@ LOCAL_ADDR=$(LOCAL_ADDR) docker-compose up -d
	@ sleep 25

import: get-cli
	@ echo "LOCAL_ADDR is "$(LOCAL_ADDR)
	@ for c in $(clusters); do \
		LOCAL_ADDR=$(LOCAL_ADDR) scalebox app create $$c/mycluster.yaml; \
	done

up:
	@ docker-compose up -d

down: clean-agents
	@ docker-compose down
	# @ docker ps -a|grep -v CONTAINER|awk '{print $$1}'| while read line; do docker rm -f $$line;done

clean-agents:
	@ # clear all agent-containers
	@ docker ps -a|grep goagent|awk '{print $$1}'| while read line; do docker rm -f $$line;done

clean:
	@ pdsh -w $(NODES) 'docker rm -f $$(docker ps -qa)' | dshbak -c
	@ echo "================"
	@ docker rm -f $$(docker ps -qa)

list:
	@ pdsh -w $(NODES) 'docker  ps -a' | dshbak -c
	@ echo "================"
	@ docker ps -a

# actuator-passwordless:
# 	@ mkdir -p $(SHARED_DIR) && cp id_rsa.pub $(SHARED_DIR)
# 	@ pdsh -w $(NODES) 'cat $(SHARED_DIR)/id_rsa.pub >> ${HOME}/.ssh/authorized_keys' | dshbak -c
# 	@ rm -f $(SHARED_DIR)/id_rsa.pub

pull-all:
	@ for i in database actuator controld agent base cli; do \
		echo $$i; \
		docker pull hub.cstcloud.cn/scalebox/$$i; \
	done

get-cli:
	@ mkdir -p ~/bin
	@ os=`uname| tr '[:upper:]' '[:lower:]'` \
	&& arch=`uname -m` \
	&& if [ "$$arch" = "x86_64" ]; then arch="amd64"; fi \
	&& cid=$$(docker create hub.cstcloud.cn/scalebox/cli) \
	&& docker cp $$cid:/usr/local/bin/scalebox-$$os-$$arch ~/bin/scalebox \
	&& docker rm -v $$cid

auto-keygen-actuator:
	@ make -C ../dockerfiles/actuator build
	@ docker tag my-actuator hub.cstcloud.cn/scalebox/actuator
	@ docker rmi my-actuator

get-pubkey:
	@ cid=`docker create hub.cstcloud.cn/scalebox/actuator` \
	&& docker cp $$cid:/root/.ssh/id_rsa.pub ${PWD} \
	&& docker rm -v $$cid

update-pubkey: auto-keygen-actuator get-pubkey

prepare: get-pubkey
	@ mkdir -p ~/.ssh ~/.scalebox/log
	@ chmod 700 ~/.ssh
	@ echo "alias app='scalebox app'" >> ~/.bash_profile
	@ echo "alias job='scalebox job'" >> ~/.bash_profile
	@ echo "alias task='scalebox task'" >> ~/.bash_profile
	@ echo "export PATH=/usr/local/bin:/usr/local/sbin:\$$PATH" >> ~/.bashrc
	@ # setup passwordless ssh from actuator to node
	@ cat id_rsa.pub >> ~/.ssh/authorized_keys
	@ chmod 600 ~/.ssh/authorized_keys
