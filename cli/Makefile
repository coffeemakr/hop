DESTDIR=dist
BINARY=$(DESTDIR)/ruck
LDFLAGS=-s -w
MAIN=app/ruck.go
SOURCES+=$(wildcard cmd/*.go) $(wildcard ../*.go)

.PHONY: all
all: $(BINARY)

$(BINARY): $(MAIN) $(SOURCES)
	go build -ldflags "$(LDFLAGS)" -o $@  $< 