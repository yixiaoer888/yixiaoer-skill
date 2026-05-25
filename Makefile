.PHONY: build test

build:
	go build -o bin/yxer.exe .

test:
	go test ./...
