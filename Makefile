all: build

build:
	@go build -o ./bin/sysinfo .

test:
	@go test -v ./...

.PHONY: build test
