DESTDIR=$(abspath ./dist)
export DESTDIR
export prefix
all: server cli


clean:
	rm -f $(BINARIES)
	$(MAKE) -C cli $@
	$(MAKE) -C server $@

server:
	$(MAKE) "DESTDIR=$(DESTDIR)" -C server

cli:
	$(MAKE) "DESTDIR=$(DESTDIR)" -C cli

%:
	$(MAKE) -C cli $@
	$(MAKE) -C server $@