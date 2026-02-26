IMAGE_NAME:=app-primes/builder:bash

build:
	docker build --network=host -t $(IMAGE_NAME) .

clean:
	docker rmi -f app-primes/builder:bash app-primes/calc:bash

run:
	docker run -it --rm -v ${PWD}:/src --entrypoint bash $(IMAGE_NAME)
	# docker run -it --rm --entrypoint bash $(IMAGE_NAME)
