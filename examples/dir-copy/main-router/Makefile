IMAGE_NAME:=dir-copy/main-router

build:
	docker build --network=host -t $(IMAGE_NAME) .

clean:
	docker rmi $(IMAGE_NAME)
