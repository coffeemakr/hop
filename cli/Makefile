DESTDIR=dist
BINARY=$(DESTDIR)/ruck

MAIN=app/ruck.go
SOURCES+=$(wildcard cmd/*.go) $(wildcard ../*.go)

.PHONY: all
all: $(BINARY)

$(BINARY): $(MAIN) $(SOURCES)
	go build -o $@  $< 