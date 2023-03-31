IMAGE_NAME:=timeout-gen

build:
	docker build --network=host -t $(IMAGE_NAME) .

clean:
	docker rmi $(IMAGE_NAME)
