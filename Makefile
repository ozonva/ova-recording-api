LOCAL_BIN:=$(CURDIR)/bin

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

bin-deps:
	go get -u google.golang.org/grpc
	go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
	go get -u github.com/golang/protobuf/proto
	go install google.golang.org/protobuf/cmd/protoc-gen-go
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.5.0
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.5.0
	go get github.com/envoyproxy/protoc-gen-validate
	go install github.com/envoyproxy/protoc-gen-validate

generate:
	protoc -I=./ --go_out=./pkg/recording --go-grpc_out=pkg/recording ./api/api.proto