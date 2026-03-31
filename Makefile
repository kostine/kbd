BINARY  := kbd
PKG     := github.com/kostine/kbd
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)

GO      := go
GOFLAGS :=

.PHONY: build run install clean fmt vet test

build:
	$(GO) build $(GOFLAGS) -o $(BINARY) .

run: build
	./$(BINARY) $(ARGS)

install:
	$(GO) install $(GOFLAGS) .

clean:
	rm -f $(BINARY)

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

test:
	$(GO) test ./... -v
