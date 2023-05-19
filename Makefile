.PHONY: build lint test

PROJECT_PACKAGE="github.com/mirkoperillo/enigma"

build:
	@go build -o ./bin/enigma ${PROJECT_PACKAGE}/cmd/enigma
	@echo "Compilation result in: ./bin/enigma"

lint:
	@golangci-lint run
test:
	@go test ${PROJECT_PACKAGE}/...

