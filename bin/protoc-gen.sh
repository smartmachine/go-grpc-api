#!/bin/bash
go build -o ./bin/generators/protoc-gen-go ./vendor/github.com/golang/protobuf/protoc-gen-go
protoc --plugin=bin/generators/protoc-gen-go --proto_path=api/proto/v1 --go_out=plugins=grpc:pkg/api/v1 todo-service.proto
