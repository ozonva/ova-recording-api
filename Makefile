.PHONY: build, run, lint

default: build

build: ./cmd/ova-recording-api/main.go
	go build -o ./build/ova-recording-api ./cmd/ova-recording-api

run:
	go run ./cmd/ova-recording-api

test:
	go test -v ./internal/utils/*.go

lint:
	$(info ******************** running lint tools ********************)
	/home/evyalyy/go/bin/golangci-lint run
