DESTDIR=dist
BINARY=$(DESTDIR)/ruckd

MAIN=app/ruckd.go
SOURCES+=$(wildcard cmd/*.go) $(wildcard ../*.go)
.PHONY: all
all: $(BINARY)

$(BINARY): $(MAIN) $(SOURCES)
	go build -o $@  $< 