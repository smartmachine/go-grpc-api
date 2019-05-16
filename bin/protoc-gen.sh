#!/bin/bash
go build -o ./bin/generators/protoc-gen-go           ./vendor/github.com/golang/protobuf/protoc-gen-go
go build -o ./bin/generators/protoc-gen-grpc-gateway ./vendor/github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
go build -o ./bin/generators/protoc-gen-swagger      ./vendor/github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
protoc --plugin=bin/generators/protoc-gen-go           --proto_path=api/proto/v1 --go_out=plugins=grpc:pkg/api/v1               todo-service.proto
protoc --plugin=bin/generators/protoc-gen-grpc-gateway --proto_path=api/proto/v1 --grpc-gateway_out=logtostderr=true:pkg/api/v1 todo-service.proto
protoc --plugin=bin/generators/protoc-gen-swagger      --proto_path=api/proto/v1 --swagger_out=logtostderr=true:api/swagger/v1  todo-service.proto

