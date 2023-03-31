IMAGE_NAME:=hub.cstcloud.cn/scalebox/ftp-copy

build:
	docker build --network=host -t $(IMAGE_NAME) .

push:
	docker push $(IMAGE_NAME)

clean:
	docker rmi $(IMAGE_NAME)
