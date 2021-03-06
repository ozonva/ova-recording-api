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
	go test ./... -v -count=1

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
$(GENERATED_API): api/api.proto api/kafka.proto
	protoc -I=./ --go_out=./pkg/recording --go-grpc_out=./pkg/recording ./api/api.proto
	protoc -I=./ --go_out=./pkg/recording ./api/kafka.proto

GENERATED_MOCKS := ./internal/repo/mock/mock_repo.go
$(GENERATED_MOCKS): ./internal/repo/repo.go
	mockgen -source ./internal/repo/repo.go -destination $(GENERATED_MOCKS)

GENERATED_MOCKS_KFK := ./internal/kafka_client/mock/mock_kafka_client.go
$(GENERATED_MOCKS_KFK): ./internal/kafka_client/kafka_client.go
	mockgen -source ./internal/kafka_client/kafka_client.go -destination $(GENERATED_MOCKS_KFK)

GENERATED_MOCKS_METRICS := ./internal/app/metrics/mock/mock_metrics.go
$(GENERATED_MOCKS_METRICS): ./internal/app/metrics/metrics.go
	mockgen -source ./internal/app/metrics/metrics.go -destination $(GENERATED_MOCKS_METRICS)
generate: $(GENERATED_API) $(GENERATED_MOCKS) $(GENERATED_MOCKS_KFK) $(GENERATED_MOCKS_METRICS)
