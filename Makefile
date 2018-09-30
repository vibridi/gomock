
clean:
	rm -r build/

deps:
	dep ensure -v

fmt:
	goimports -w error/ helper/ parser/ writer/ main.go

build:
	go build -o build/gomock

test:
	go test -v ./...
