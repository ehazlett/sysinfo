all: build

build:
	@go build -o ./bin/sysinfo .

.PHONY: build
