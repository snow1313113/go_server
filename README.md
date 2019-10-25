# go_server
go rpc server

## how to gen pb.go
build protoc first
./get.sh github.com/golang/protobuf/protoc-gen-go
cp src/plugin_link/link_pepper.go src/github.com/golang/protobuf/protoc-gen-go
./build.sh github.com/golang/protobuf/protoc-gen-go
./protoc --go_out=plugins=pepper:. --plugin=protoc-gen-go *.proto

## how to start simple_server
./build.sh simple_server
./bin/simple_server

## how to start simple_client
./build.sh simple_client
./bin/simple_client


