IMAGE_NAME:=hub.cstcloud.cn/scalebox/cluster-file-copy

build:
	DOCKER_BUILDKIT=1 docker build --network=host -t $(IMAGE_NAME) .

push:
	docker push $(IMAGE_NAME)

clean:
	docker rmi $(IMAGE_NAME)

dist:
	docker save $(IMAGE_NAME) | zstdmt | pv | ssh c0 'zstd -d | docker load'
