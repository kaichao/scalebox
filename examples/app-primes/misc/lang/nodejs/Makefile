IMAGE_NAME:=app-primes/builder:nodejs

build:
	docker build --network=host -t $(IMAGE_NAME) .

clean:
	docker rmi -f app-primes/builder:nodejs app-primes/calc:nodejs
	docker rmi -f node:25-bookworm-slim
