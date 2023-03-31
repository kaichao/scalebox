IMAGE_NAME:=hub.cstcloud.cn/scalebox/rsync-copy
IMAGE_FILE:=/data/tmp/my.img

build:
	docker build --network=host -t $(IMAGE_NAME) .

push:
	docker push $(IMAGE_NAME)

clean:
	docker rmi $(IMAGE_NAME)

run:
	docker run -it --rm --entrypoint bash $(IMAGE_NAME)

dist:
	docker save $(IMAGE_NAME) > $(IMAGE_FILE)
	pdsh -w n[0-3] "docker load < $(IMAGE_FILE)" | dshbak -c
	# for i in {0..3}; do echo $$i; ssh n$$i docker load < $(IMAGE_FILE); done 
	rm -f $(IMAGE_FILE)
