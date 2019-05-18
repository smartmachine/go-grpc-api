#!/bin/bash
go get -d github.com/infobloxopen/protoc-gen-gorm
CURRENTPWD=$(pwd)
cd ${GOPATH}/src/github.com/infobloxopen/protoc-gen-gorm
git checkout v0.16.0
dep ensure -v
make install
cd ${CURRENTPWD}
