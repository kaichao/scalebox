IMAGE_NAME:=app-primes/builder:postgres
SLAVE_IMAGE:=app-primes/postgres-operator

build:
	docker build --network=host -f Dockerfile.master -t $(IMAGE_NAME) .
	docker build --network=host -f Dockerfile.slave -t $(SLAVE_IMAGE) .

run:
	# docker run -it --rm -v ${PWD}:/src --entrypoint bash $(IMAGE_NAME)
	# docker run -it --rm -p 5432:5432 --entrypoint bash $(IMAGE_NAME)
	docker run -v ${PWD}:/src -e POSTGRES_PASSWORD=changeme -d sameersbn/postgresql:12-20200524
