# Default workspace.
WORKSPACE ?= $(shell pwd)
export WORKSPACE

.PHONY: build run go lint test

# Execution.
go: build run 
build:
	go build -o ${WORKSPACE}/bin/walletapp ${WORKSPACE}/cmd/main.go
run:
	go run ${WORKSPACE}/bin/walletapp

# Utilitary.
lint:
	golangci-lint run ./...
test:
	go test -v -race ./...