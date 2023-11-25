IMAGE_NAME:=dir-copy/message-router

build:
	docker build --network=host -t $(IMAGE_NAME) .

clean:
	docker rmi $(IMAGE_NAME)
