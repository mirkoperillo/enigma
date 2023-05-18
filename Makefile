.PHONY: build test

PROJECT_PACKAGE="github.com/mirkoperillo/enigma"

build:
	@go build -o ./bin/enigma ${PROJECT_PACKAGE}/cmd/enigma
	@echo "Compilation result in: ./bin/enigma"

test:
	@go test ${PROJECT_PACKAGE}/...

