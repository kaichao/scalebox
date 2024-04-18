IMAGE_NAME:=hub.cstcloud.cn/scalebox/rsync-copy

build:
	DOCKER_BUILDKIT=1 docker build --network=host -t $(IMAGE_NAME) .

push:
	docker push $(IMAGE_NAME)

pull:
	docker pull $(IMAGE_NAME)

clean:
	docker rmi $(IMAGE_NAME)
