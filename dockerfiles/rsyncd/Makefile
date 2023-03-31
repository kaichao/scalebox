IMAGE_NAME:=hub.cstcloud.cn/scalebox/rsyncd

build:
	docker build --network=host -t $(IMAGE_NAME) .

push:
	docker push $(IMAGE_NAME)

clean:
	docker rmi $(IMAGE_NAME)

run-rsync-mode:
	docker run -d --rm --network host -e TZ=CST-8 -e RSYNC_PORT=873 \
		-v /etc/localtime:/etc/localtime:ro -v /var/log:/var/log -v /:/local $(IMAGE_NAME) 

run-rsync-over-ssh:
	docker run -d --rm --network host -e TZ=CST-8 -e TRANSPORT_TYPE=ssh -e SSH_PORT=2222 \
		-v /etc/localtime:/etc/localtime:ro -v /var/log:/var/log -v /:/local $(IMAGE_NAME) 
