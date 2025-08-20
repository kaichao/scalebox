IMAGE_NAME:=hub.cstcloud.cn/scalebox/file-copy:latest

build:
	DOCKER_BUILDKIT=1 docker build --network=host -t $(IMAGE_NAME) .

build-singularity:
	@mkdir -p ~/singularity/scalebox/
	@date
	@docker save $(IMAGE_NAME) -o /tmp/file-copy.tar
	@singularity build -F ~/singularity/scalebox/file-copy.sif docker-archive:///tmp/file-copy.tar
	@rm -f /tmp/file-copy.tar
	@date

push:
	docker push $(IMAGE_NAME)

pull:
	docker pull $(IMAGE_NAME)

clean:
	docker rmi $(IMAGE_NAME)

run: build
	docker run -it --rm --entrypoint bash $(IMAGE_NAME)
