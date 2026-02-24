IMAGE_NAME:=app-primes/calc:python

build:
	docker build --network=host -t $(IMAGE_NAME) .

clean:
	docker rmi $(IMAGE_NAME)

dist: build
	# pdsh -w n[0-3] "docker load < /data/tmp/my.img" | dshbak -c
	for i in {0..3}; do echo $$i;docker save $(IMAGE_NAME) | pv | ssh n$$i 'docker load';done
