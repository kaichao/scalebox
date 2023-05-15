IMAGE_NAME:=my-actuator

build:
	docker build --network=host -t $(IMAGE_NAME) .

clean:
	docker rmi $(IMAGE_NAME)
