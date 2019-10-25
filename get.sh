#!/bin/bash
if [ $# -lt 1 ]; then
    echo "USAGE: $0 url"
    echo "e.g.: $0 github.com/golang/protobuf/protoc-gen-go"
    exit -1
fi

cd `dirname $0`
CURRENT_DIR=`pwd`
export GOPATH=${CURRENT_DIR}
go get -d -u $1
