#!/bin/bash
if [ $# -lt 1 ]; then
    echo "USAGE: $0 sub_dir"
    echo "e.g.: $0 test"
    exit -1
fi

cd `dirname $0`
CURRENT_DIR=`pwd`
#export GOPATH=${CURRENT_DIR}
cd src/$1
go build -o ../../../bin/ .
#go install .
