IMAGE_NAME:=app-primes/builder:php

build:
	docker build --network=host -t $(IMAGE_NAME) .

run:
	docker run -it --rm -v ${PWD}:/src $(IMAGE_NAME) sh
