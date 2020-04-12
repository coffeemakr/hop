prefix=/usr/local
BUILDIR=build
BINARY=$(BUILDIR)/ruck
LDFLAGS=-s -w
MAIN=app/ruck.go
SOURCES+=$(wildcard cmd/*.go) $(wildcard ../*.go)

.PHONY: all
all: $(BINARY)

$(BINARY): $(MAIN) $(SOURCES)
	go build -ldflags "$(LDFLAGS)" -o $@  $< 

.PHONY: install
install: $(BINARY)
	install -D $(BINARY) $(DESTDIR)$(prefix)/bin/ruck