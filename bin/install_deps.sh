#!/bin/bash

# Check for the presence of git.  Exit if not available.
if ! [ -x "$(command -v git)" ]; then
  echo 'Error: git is not installed. Please install it on this system.' >&2
  exit 1
fi

# Check for the presence of golang.  Exit if not available.
if ! [ -x "$(command -v go)" ]; then
  echo 'Error: go is not installed. Please install it on this system.' >&2
  exit 1
fi

# Check for the presence of Google Protocol Buffers.  Exit if not available.
if ! [ -x "$(command -v protoc)" ]; then
  echo 'Error: Google Protocol Buffers is not installed. Please install it on this system.' >&2
  exit 1
fi

# Check for the presence of curl.  Exit if not available.
if ! [ -x "$(command -v curl)" ]; then
  echo 'Error: curl is not installed. Please install it on this system.' >&2
  exit 1
fi

# Check if dep is installed, if not, get the latest version
if ! [ -x "$(command -v dep)" ]; then
  echo 'Warning: dep is not installed.  Installing it ...' >&2
  curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh >&2
fi

# Check if protoc-gen-go is installed, if not, install it.
if ! [ -x "$(command -v protoc-gen-go)" ]; then
  echo 'Warning: protoc-gen-go is not installed.  Installing it ...' >&2
  go get -u github.com/golang/protobuf/protoc-gen-go >&2
fi

# Check if protoc-gen-swagger is installed, if not, install it.
if ! [ -x "$(command -v protoc-gen-swagger)" ]; then
  echo 'Warning: protoc-gen-swagger is not installed.  Installing it ...' >&2
  go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
fi

# Check if protoc-gen-grpc-gateway is installed, if not, install it.
if ! [ -x "$(command -v protoc-gen-grpc-gateway)" ]; then
  echo 'Warning: protoc-gen-grpc-gateway is not installed.  Installing it ...' >&2
  go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
fi

# Check if protoc-gen-gorm is installed, if not, install it.
if ! [ -x "$(command -v protoc-gen-gorm)" ]; then
  # protoc-gen-gorm currently has a nasty dependency bug, can't be go getted.
  echo 'Warning: protoc-gen-gorm is not installed.  Installing it ...' >&2
  go get -d github.com/infobloxopen/protoc-gen-gorm
  CURRENTPWD=$(pwd)
  cd ${GOPATH}/src/github.com/infobloxopen/protoc-gen-gorm
  git checkout v0.16.0
  dep ensure -v
  make install
  cd ${CURRENTPWD}
fi
