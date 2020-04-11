DESTDIR=dist
BINARY=$(DESTDIR)/ruckd
LDFLAGS=-s -w
MAIN=app/ruckd.go
SOURCES+=$(wildcard cmd/*.go) $(wildcard ../*.go)
.PHONY: all
all: $(BINARY)

$(BINARY): $(MAIN) $(SOURCES)
	go build -ldflags "$(LDFLAGS)" -o $@  $< 