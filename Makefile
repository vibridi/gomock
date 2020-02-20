.PHONY: fmt deps test clean build

VERSION = $(shell git describe --tags)
GOVERSION = $(shell go version | cut -c 12-)
BUILD_LDFLAGS = "\
          -X \"github.com/vibridi/gomock/version.VERSION=$(VERSION)\" \
          -X \"github.com/vibridi/gomock/version.GOVERSION=$(GOVERSION)\""

clean:
	rm -rf ./build/

deps:
	go get && go mod tidy

fmt:
	goimports -w ./_example/ error/ helper/ parser/ version/ writer/ main.go

build: clean
	go build -ldflags=$(BUILD_LDFLAGS) -o ./build/gomock *.go

example: build
	./build/gomock -f _example/_example.go

example-qualify: build
	./build/gomock -f _example/_qualify.go -q

example-export: build
	./build/gomock -f _example/_example.go -x

example-compose: build
	./build/gomock -f _example/_composition.go

install: build
	mv ./build/gomock $(GOPATH)/bin/

test:
	go test -v -cover ./...
