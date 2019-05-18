#!/bin/bash

declare -a SYSTEM_DEPS
declare -a GO_DEP_NAME
declare -a GO_DEP_INSTALL

SYSTEM_DEPS=(git go protoc curl)
GO_DEP_NAME=(dep protoc-gen-go protoc-gen-swagger protoc-gen-grpc-gateway protoc-gen-gorm)
GO_DEP_INSTALL=(
"curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh"
"go get -u github.com/golang/protobuf/protoc-gen-go"
"go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger"
"go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway"
"bin/install_protoc-gen-gorm.sh"
)

check_system_dep() {
  if ! [ -x "$(command -v $1)" ]; then
    echo "Error: $1 is not installed. Please install it on this system."
    exit 1
  fi
}

check_go_dep() {
  if ! [ -x "$(command -v $1)" ]; then
    echo "Warning: $1 is not installed.  Installing it ..."
    eval $2
  fi
}

for i in "${SYSTEM_DEPS[@]}" ; do
  check_system_dep $i
done

for i in "${!GO_DEP_NAME[@]}" ; do
  check_go_dep "${GO_DEP_NAME[$i]}" "${GO_DEP_INSTALL[$i]}"
done

