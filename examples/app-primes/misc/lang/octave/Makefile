IMAGE_NAME:=app-primes/builder:octave

build:
	docker build --network=host -t $(IMAGE_NAME) .

run:
	docker run -it --rm -v ${PWD}:/src --entrypoint bash $(IMAGE_NAME)

clean:
	docker rmi $(IMAGE_NAME)
