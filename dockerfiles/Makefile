dirs:=filelist ftp-copy rsync-agent rsyncd rsync-copy

build:
	@for dir in $(dirs); do \
		$(MAKE) -C $$dir build; \
	done

push:
	@for dir in $(dirs); do \
		$(MAKE) -C $$dir push; \
	done

clean:
	@for dir in $(dirs); do \
		$(MAKE) -C $$dir clean; \
	done
