IMAGE_NAME:=dir-copy/message-router

build:
	DOCKER_BUILDKIT=1 docker build --network=host -t $(IMAGE_NAME) .

clean:
	docker rmi $(IMAGE_NAME)
