# go_server
go rpc server

## prepare
1. build protoc first

2. ./get.sh github.com/golang/protobuf/protoc-gen-go

3. cp src/pepper/link_pepper src/github.com/golang/protobuf/protoc-gen-go/link_pepper.go

4. ./build.sh github.com/golang/protobuf/protoc-gen-go

## how to gen pb.go
1. Set your protobuf dir in *gen_pb_go.sh*

2. ./gen_pb_go.sh

## how to start simple_server
1. ./build.sh server/simple_server

2. ./bin/simple_server

## how to start simple_client
1. ./build.sh test/simple_client

2. ./bin/simple_client
