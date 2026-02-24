IMAGE_NAME:=app-primes/builder:r

build:
	docker build --network=host -t $(IMAGE_NAME) .

run:
	docker run -it --rm -v ${PWD}:/src --entrypoint bash $(IMAGE_NAME)

clean:
	docker kill $(CONTAINER_NAME) || true
	docker rmi $(IMAGE_NAME)
