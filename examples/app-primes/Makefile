DIRS ?= calc

build:
	@for d in $(DIRS); do make -C $$d build; done

clean:
	@for d in $(DIRS); do make -C $$d clean; done

dist:
	for i in {0..3}; do echo $$i;docker save hub.cstcloud.cn/scalebox/agent | pv | ssh n$$i 'docker load';done

dist-all: dist
	make -C calc dist

sync:
	rsync -av . --del bio:/shared/app-primes
	
