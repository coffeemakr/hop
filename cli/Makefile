DESTDIR=dist
BINARY=$(DESTDIR)/hophop

MAIN=app/hophop.go
SOURCES+=$(wildcard cmd/*.go) $(wildcard ../*.go)

.PHONY: all
all: $(BINARY)

$(BINARY): $(MAIN) $(SOURCES)
	go build -o $@  $< 