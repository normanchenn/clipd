.PHONY: build run test tidy clean install help

build:
	go build -o ./bin/clipd ./cmd/clipd

run:
	./bin/clipd

test:
	go test ./internal/... ./cmd/... -v

tidy:
	go mod tidy

clean:
	rm -rf ./bin/*

install:
	./scripts/install.sh

help:
	@echo "TODO"

