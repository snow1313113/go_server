#!/bin/bash
if [ $# -lt 1 ]; then
    echo "USAGE: $0 sub_dir"
    echo "e.g.: $0 test"
    exit -1
fi

export GO111MODULE=on
export GOPROXY=https://goproxy.io
cd `dirname $0`
cd src/$1
go build -o ../../../bin/ .
