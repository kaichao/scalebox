# LANG_TAGS ?= c fortran python java golang rust nodejs r octave bash php postgres julia tcl
LANGS ?= c fortran python java golang rust nodejs r php octave bash postgres
# LANG_TAGS ?= r php octave bash postgres
# LANG_TAGS ?= r octave bash postgres
# LANG_TAGS ?= c
SHARED_DIR ?= /gfsdata

ifeq ($(TAG), java)
  BASE_IMAGE := openjdk:11-jre-bullseye
else ifeq ($(TAG), python)
  BASE_IMAGE := python:3.14-slim
else ifeq ($(TAG), nodejs)
  BASE_IMAGE := node:25-bookworm-slim
else ifeq ($(TAG), octave)
  BASE_IMAGE := mtmiller/octave:7.0.0
else ifeq ($(TAG), r)
  BASE_IMAGE := r-base:4.5.2
else ifeq ($(TAG), php)
  BASE_IMAGE := php:8-cli
else ifeq ($(TAG), postgres)
  BASE_IMAGE := app-primes/builder:postgres
else ifeq ($(TAG), tcl)
  BASE_IMAGE := efrecon/tcl
else
  BASE_IMAGE := debian:13-slim
endif

APP_IMAGE := app-primes/calc:$(TAG)
BUILDER_IMAGE := app-primes/builder:$(TAG)

build-all:
	@for tag in $(LANGS); do TAG=$$tag make build; done

push-all:
	@for tag in $(LANGS); do TAG=$$tag make push; done

run-all:
	@for tag in $(LANGS); do echo LANG=$$tag; TAG=$$tag make run; done

clean-all:
	@for tag in $(LANGS); do TAG=$$tag make clean; done

dist-all:
	@for tag in $(LANGS); do TAG=$$tag make dist; done

dist:
	docker save app-primes/$(TAG)-calc > $(SHARED_DIR)/tmp/my.img
	pdsh -w r0[1-9],r1[0-9] -l root "docker load < $(SHARED_DIR)/tmp/my.img" | dshbak -c
	rm -f $(SHARED_DIR)/tmp/my.img

dist-operator:
	docker save app-primes/postgres-operator > $(SHARED_DIR)/tmp/my.img
	pdsh -w r0[1-9],r1[0-9] -l root "docker load < $(SHARED_DIR)/tmp/my.img" | dshbak -c
	rm -f $(SHARED_DIR)/tmp/my.img


build:
	make -C $(TAG) build
	docker build --network=host \
		--build-arg BASE_IMAGE=$(BASE_IMAGE) \
		--build-arg BUILDER_IMAGE=$(BUILDER_IMAGE) \
		--tag $(APP_IMAGE) .

push:
	docker push $(APP_IMAGE)

run:
	@docker run -it --rm $(BUILDER_IMAGE) /src/primes 1 10000

clean:
	docker rmi $(APP_IMAGE)
	docker rmi $(BUILDER_IMAGE)

