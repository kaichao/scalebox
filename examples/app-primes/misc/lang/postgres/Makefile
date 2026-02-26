IMAGE_NAME:=app-primes/builder:postgres
SLAVE_IMAGE:=app-primes/postgres-operator

build:
	docker build --network=host -f Dockerfile.master -t $(IMAGE_NAME) .
	docker build --network=host -f Dockerfile.slave -t $(SLAVE_IMAGE) .

clean:
	docker rmi -f app-primes/builder:postgres app-primes/calc:postgres app-primes/postgres-operator
	docker rmi -f postgres:18 kaichao/docker-cli
