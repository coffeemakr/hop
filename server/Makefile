DESTDIR=dist
BINARY=$(DESTDIR)/server

MAIN=app/server.go
SOURCES+=$(wildcard cmd/*.go) $(wildcard ../*.go)
.PHONY: all
all: $(BINARY)

$(BINARY): $(MAIN) $(SOURCES)
	go build -o $@  $< 