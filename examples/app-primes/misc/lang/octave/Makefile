IMAGE_NAME:=app-primes/builder:octave

build:
	docker build --network=host -t $(IMAGE_NAME) .

clean:
	docker rmi -f app-primes/builder:octave app-primes/calc:octave
	docker rmi -f mtmiller/octave:7.0.0
