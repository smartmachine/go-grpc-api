VERSION := $(shell git describe --tags --dirty)
BUILD := $(shell git rev-parse --short HEAD)
GOFILES := $(wildcard *.go)
LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

.PHONY: all check ensure api test dep build clean clean-api veryclean help

all: check ensure api test build

check: ## Install all supporting dependencies like generators and dep.
	$(info Checking system for required build/compile-time dependencies)
	@bin/install_deps.sh

Gopkg.lock: Gopkg.toml $(GOFILES)
	$(info Rebuilding Gopkg.lock)
	@dep ensure -update

ensure: Gopkg.lock ## Ensure vendor directory is up to date

pkg/api/v1/todo-service.pb.go: api/proto/v1/todo-service.proto
	$(info ... Generating Protobuffer Go files)
	@protoc --proto_path=third_party --proto_path=api/proto/v1 --go_out=plugins=grpc:pkg/api/v1 todo-service.proto

api/swagger/v1/todo-service.swagger.json: api/proto/v1/todo-service.proto
	$(info ... Generating Swagger Documentation)
	@protoc --proto_path=third_party --proto_path=api/proto/v1 --swagger_out=logtostderr=true:api/swagger/v1 todo-service.proto

pkg/api/v1/todo-service.pb.gw.go: api/proto/v1/todo-service.proto
	$(info ... Generating GRPC Gateway [REST] proxy)
	@protoc --proto_path=third_party --proto_path=api/proto/v1 --grpc-gateway_out=logtostderr=true:pkg/api/v1 todo-service.proto

pkg/api/v1/todo-service.pb.gorm.go: api/proto/v1/todo-service.proto
	$(info ... Generating GORM Protobuffer->ORM structures)
	@protoc --proto_path=third_party --proto_path=api/proto/v1 --gorm_out=logtostderr=true:pkg/api/v1 todo-service.proto

api: pkg/api/v1/todo-service.pb.go api/swagger/v1/todo-service.swagger.json pkg/api/v1/todo-service.pb.gw.go pkg/api/v1/todo-service.pb.gorm.go ## Auto-generate grpc go sources

test: ## Run unit tests
	$(info Running unit tests ...)
	@go test ./pkg/service/v1

dep: ## Make sure all dependencies are up to date
	$(info Checking vendor directory)
	@dep ensure

server: $(GOFILES)
	$(info Compiling server)
	@go build -v $(LDFLAGS) -o server ./cmd/server

client: $(GOFILES)
	$(info Compiling client)
	@go build -v $(LDFLAGS) -o client ./cmd/client

client-rest: $(GOFILES)
	$(info Compiling REST client)
	@go build -v $(LDFLAGS) -o client-rest ./cmd/client-rest

build: dep server client client-rest ## Build all binary artifacts

clean: ## Clean all build artifacts
	$(info Cleaning all build artifacts)
	@rm -rf server client client-rest
	@go clean

clean-api: ## Remove all generated code and files.  Regenerate with api target.
	$(info Removing all generated code and files)
	@rm -rfv pkg/api/v1/todo-service.pb.go api/swagger/v1/todo-service.swagger.json pkg/api/v1/todo-service.pb.gw.go pkg/api/v1/todo-service.pb.gorm.go

veryclean: clean clean-api ## Clean all caches and generated objects
	@go clean -cache -testcache -modcache

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

