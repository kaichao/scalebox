dirs:=cron list-dir ftp-copy rsync-copy rsyncd

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
