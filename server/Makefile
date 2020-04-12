BUILDIR=build
prefix=/usr/local
BINARY=$(BUILDIR)/ruckd
LDFLAGS=-s -w
MAIN=app/ruckd.go
SOURCES+=$(wildcard cmd/*.go) $(wildcard ../*.go)
.PHONY: all
all: $(BINARY)

$(BINARY): $(MAIN) $(SOURCES)
	go build -ldflags "$(LDFLAGS)" -o $@  $< 

.PHONY: install
install: $(BINARY)
	install -D $(BINARY) $(DESTDIR)$(prefix)/bin/ruckd
