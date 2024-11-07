# Setting up workspace
WORKSPACE ?= $(shell pwd)
export WORKSPACE

-include ${WORKSPACE}/config.env
export $(shell sed 's/=.*//' config.env)

.PHONY: build run go lint test

# Execution.
go: build run 
build:
	go build -o ${WORKSPACE}/bin/walletapp ${WORKSPACE}/cmd/main.go
run:
	${WORKSPACE}/bin/walletapp

# Utilitary.
lint:
	golangci-lint config --config ./.golangci.yaml
	golangci-lint run ./...
test:
	go test -v -race ./...