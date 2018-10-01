.PHONY: fmt deps test clean build

VERSION = $(shell git describe --tags)
GOVERSION = $(shell go version | cut -c 12-)
BUILD_LDFLAGS = "\
          -X \"github.com/vibridi/gomock/version.VERSION=$(VERSION)\" \
          -X \"github.com/vibridi/gomock/version.GOVERSION=$(GOVERSION)\""

clean:
	rm -rf build/

deps:
	dep ensure -v

fmt:
	goimports -w _example/ error/ helper/ parser/ version/ writer/ main.go

build: clean
	go build -ldflags=$(BUILD_LDFLAGS) -o build/gomock *.go

test:
	go test -v -cover ./...
