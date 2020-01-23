default: build

GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)

build:
	GO111MODULE=on
	go build -v -ldflags="-s -w" -o .target/cleanup main.go

.PHONY: clean
clean:
	rm -rf .target

.PHONY: test
test:
	go test -v ./...
