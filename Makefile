VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GO ?= go
LDFLAGS := -ldflags "-s -w -X github.com/hoalong/lume-fleet/cmd.Version=$(VERSION)"

.PHONY: build install clean

build:
	$(GO) build $(LDFLAGS) -o lume-fleet .

install:
	$(GO) install $(LDFLAGS) .

clean:
	rm -f lume-fleet
