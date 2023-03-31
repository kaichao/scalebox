BASE_IMAGE:=hub.cstcloud.cn/scalebox/rsync-base
AGENT_IMAGE:=hub.cstcloud.cn/scalebox/rsync-agent

build:
	docker build --network=host --target agent-base -t $(BASE_IMAGE) .
	docker build --network=host -t $(AGENT_IMAGE) .

push:
	docker push $(BASE_IMAGE)
	docker push $(AGENT_IMAGE)

clean:
	docker rmi $(BASE_IMAGE)
	docker rmi $(AGENT_IMAGE)

