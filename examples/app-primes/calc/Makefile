IMAGE_NAME:=app-primes/calc

build:
	docker build --network=host -t $(IMAGE_NAME) .

push:
	docker push $(IMAGE_NAME)

clean:
	docker rmi $(IMAGE_NAME)

dist:
	# pdsh -w n[0-3] "docker load < /data/tmp/my.img" | dshbak -c
	for i in {0..3}; do echo $$i;docker save $(IMAGE_NAME) | pv | ssh n$$i 'docker load';done
