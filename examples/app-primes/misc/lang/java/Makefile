LANG_TAG:=java
IMAGE_NAME:=app-primes/builder:java

build:
	docker build --network=host -t $(IMAGE_NAME) .

run:
	docker run -it --rm $(IMAGE_NAME) bash
	# docker run -it -e NUM_PARALLEL=4 $(IMAGE_NAME) bash

exec:
	docker exec -it $(CONTAINER_NAME) bash

build-graalvm:
	docker build --network=host -f Dockerfile.graalvm -t app-primes/graalvm .

run-graalvm:
	docker run -it --rm --entrypoint bash app-primes/graalvm
