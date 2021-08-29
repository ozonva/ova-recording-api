GO_BIN=$(shell go env GOPATH)/bin
PATH:=${PATH}:${GO_BIN}

.PHONY: all
all: generate lint test run

build: ./cmd/ova-recording-api/main.go
	go build -o ./build/ova-recording-api ./cmd/ova-recording-api

.PHONY: run
run: ./cmd/ova-recording-api/main.go
	go run ./cmd/ova-recording-api --config config/config.yml

.PHONY: test
test:
	go test -v ./internal/utils/*.go
	go test -v ./internal/app/recording/*.go
	go test -v ./internal/flusher/*.go
	go test -v ./internal/repo/*.go
	go test -v ./internal/saver/*.go

.PHONY: lint
lint:
	$(info ******************** running lint tools ********************)
	golangci-lint run

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


GENERATED_API := $(wildcard ./pkg/recording/api/*.go)
$(GENERATED_API): api/api.proto
	protoc -I=./ --go_out=./pkg/recording --go-grpc_out=./pkg/recording ./api/api.proto

GENERATED_MOCKS := ./internal/repo/mock/mock_repo.go
$(GENERATED_MOCKS): ./internal/repo/repo.go
	mockgen -source ./internal/repo/repo.go -destination $(GENERATED_MOCKS)

generate: $(GENERATED_API) $(GENERATED_MOCKS)
