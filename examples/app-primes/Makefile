DIRS ?= scatter calc gather

build:
	@for d in $(DIRS); do make -C $$d build; done

push:
	@for d in $(DIRS); do make -C $$d push; done

clean:
	@for d in $(DIRS); do make -C $$d clean; done

dist:
	make -C calc dist
