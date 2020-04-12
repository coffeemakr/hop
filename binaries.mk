BUILDIR=build
prefix=/usr/local

ARM_BUILDDIR=$(BUILDIR)/armv7
ARM_DISTDIR=$(ARM_BUILDDIR)/dist
ARM_BINARY=$(addprefix $(ARM_BUILDDIR)/,$(BINARY))

AMD64_BUILDDIR=$(BUILDIR)/amd64
AMD64_DISTDIR=$(AMD64_BUILDDIR)/dist
AMD64_BINARY=$(addprefix $(AMD64_BUILDDIR)/,$(BINARY))

INSTALLED=$(DESTDIR)$(prefix)/bin/ruckd

BINARIES=$(ARM_BINARY) $(AMD64_BINARY) 
LDFLAGS=-s -w


SOURCES+=$(wildcard cmd/*.go) $(wildcard *.go) $(wildcard ../*.go)
GO=go build -ldflags "$(LDFLAGS)"
.PHONY: all
all: $(BINARIES)

$(BINARY): $(MAIN) $(SOURCES)
	$(GO) -o $@  $< 

$(ARM_BINARY): $(MAIN) $(SOURCES)
	CGO_ENABLED=0 GOOS=linux GOARM=7 GOARCH=arm $(GO) -o $@  $< 

$(AMD64_BINARY): $(MAIN) $(SOURCES)
	GOARCH=amd64 $(GO) -o $@  $< 

.PHONY: clean
clean:
	rm -f $(BINARIES)

.PHONY: install
install: $(BINARY)
	install -D $< $(DESTDIR)$(prefix)/bin/ruckd

