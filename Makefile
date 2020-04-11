DESTDIR=$(abspath ./dist)
.PHONY: all server cli

all: cli server 

server:
	$(MAKE) "DESTDIR=$(DESTDIR)" -C server

cli:
	$(MAKE) "DESTDIR=$(DESTDIR)" -C cli