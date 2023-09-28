IMAGE_BASE:=hub.cstcloud.cn/scalebox/rsyncd
IMAGE_AGENT:=hub.cstcloud.cn/scalebox/rsyncd-agent

build:
	docker build --network=host --target base -t $(IMAGE_BASE) .
	docker build --network=host -t $(IMAGE_AGENT) .
	
push:
	docker push $(IMAGE_BASE)
	docker push $(IMAGE_AGENT)

clean:
	docker rmi $(IMAGE_BASE)
	docker rmi $(IMAGE_AGENT)

run-rsync-mode:
	docker run -d --rm --network host -e TZ=CST-8 -e RSYNC_PORT=873 \
		-v /etc/localtime:/etc/localtime:ro -v /var/log:/var/log -v /:/local $(IMAGE_BASE) 

run-rsync-over-ssh:
	docker run -d --rm --network host -e TZ=CST-8 -e TRANSPORT_TYPE=ssh -e SSH_PORT=2222 \
		-v /etc/localtime:/etc/localtime:ro -v /var/log:/var/log -v /:/local $(IMAGE_BASE) 
