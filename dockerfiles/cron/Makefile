IMAGE_NAME:=hub.cstcloud.cn/scalebox/cron

build:
	DOCKER_BUILDKIT=1 docker build --network=host -t $(IMAGE_NAME) .

push:
	docker push $(IMAGE_NAME)

clean:
	docker rmi $(IMAGE_NAME)

run:
	scalebox app create
