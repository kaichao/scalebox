IMAGE_NAME:=app-primes/builder:scala

build:
	docker build --network=host -t $(IMAGE_NAME) .

clean:
	docker rmi -f app-primes/builder:scala app-primes/calc:scala
	docker rmi -f sbtscala/scala-sbt:eclipse-temurin-17.0.15_6_1.12.4_2.12.21
