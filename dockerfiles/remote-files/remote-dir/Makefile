IMAGE_NAME:=hub.cstcloud.cn/scalebox/remote-dir

build:
	DOCKER_BUILDKIT=1 docker build --network=host \
		-f Dockerfile -t $(IMAGE_NAME) ..

push:
	docker push $(IMAGE_NAME)

pull:
	docker pull $(IMAGE_NAME)

clean:
	docker rmi $(IMAGE_NAME)

run: build
	docker run -it --rm --entrypoint bash $(IMAGE_NAME)
