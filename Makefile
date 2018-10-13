
VERSION = $(shell git describe --tags)
GOVERSION=$(shell go version | cut -c 12-)
BUILD_LDFLAGS = "\
          -X \"golang.nulab-inc.com/cacoo/service/account/version.VERSION=$(VERSION)\" \
          -X \"golang.nulab-inc.com/cacoo/service/account/version.GOVERSION=$(GOVERSION)\""

clean:
	rm -rf build/

deps:
	dep ensure -v

fmt:
	goimports -w _example/ error/ helper/ parser/ version/ writer/ main.go

build: clean
	go build -ldflags=$(BUILD_LDFLAGS) -o build/gomock *.go

example: build
	./build/gomock -f _example/_example.go

test:
	go test -v -cover ./...
