# go_server
go rpc server

## prepare
1. build protoc first

2. ./get.sh github.com/golang/protobuf/protoc-gen-go

3. cp src/plugin_link/link_pepper.go src/github.com/golang/protobuf/protoc-gen-go

4. ./build.sh github.com/golang/protobuf/protoc-gen-go

## how to gen pb.go
1. Set your protobuf dir in src/protocol/gen_pb_go.sh

2. cd src/protocol

3. ./gen_pb_go.sh

## how to start simple_server
1. ./build.sh simple_server

2. ./bin/simple_server

## how to start simple_client
1. ./build.sh simple_client

2. ./bin/simple_client
