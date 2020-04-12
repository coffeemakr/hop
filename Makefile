DESTDIR=$(abspath ./dist)
export DESTDIR
export prefix
all: server cli

clean:
	rm -f $(BINARIES)
	$(MAKE) -C cli $@
	$(MAKE) -C server $@


%:
	$(MAKE) -C cli $@
	$(MAKE) -C server $@