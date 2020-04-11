DESTDIR=$(abspath ./dist)
SERVER=$(DESTDIR)/ruckd
CLI=$(DESTDIR)/ruck
BINARIES=$(SERVER) $(CLI)
.PHONY: all server cli

all: server cli


clean:
	rm -f $(BINARIES)

server:
	$(MAKE) "DESTDIR=$(DESTDIR)" -C server

cli:
	$(MAKE) "DESTDIR=$(DESTDIR)" -C cli
