.PHONY: build test

build:
	go build -o bin/yxer.exe .

test:
	powershell -ExecutionPolicy Bypass -File scripts/run-checks.ps1
