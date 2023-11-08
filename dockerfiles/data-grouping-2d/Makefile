IMAGE_NAME:=hub.cstcloud.cn/scalebox/data-grouping-2d

build:
	DOCKER_BUILDKIT=1  docker build --network=host -t $(IMAGE_NAME) .

push:
	docker push $(IMAGE_NAME)

pull:
	docker pull $(IMAGE_NAME)

clean:
	docker rmi $(IMAGE_NAME)

dist:
	docker save $(IMAGE_NAME) | zstdmt | pv | ssh c0 'zstd -d | docker load'
