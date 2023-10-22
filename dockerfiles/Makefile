dirs:=cron list-dir ftp-copy rsync-copy rsyncd data-grouping-2d cluster-file-copy cluster-dir-list

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
