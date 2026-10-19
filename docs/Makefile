IMAGE_NAME:=kaichao/read-the-docs

build-cn:
	@ docker run --rm -v $(CURDIR):/docs -w /docs/cn $(IMAGE_NAME) make html

build-en:
	@ docker run --rm -v $(CURDIR):/docs -w /docs/en $(IMAGE_NAME) make html

clean:
	rm -rf en/build/* cn/build/*
