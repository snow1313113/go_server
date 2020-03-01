# P.S PROTOBUF_DIR=~/protobuf/src
PROTOBUF_DIR=your protobuf dir
PROTOC_TOOL=${PROTOBUF_DIR}/protoc
${PROTOC_TOOL} --plugin=bin/protoc-gen-go --go_out=plugins=pepper:src/protocol/ --proto_path=${PROTOBUF_DIR}:src/protocol/ src/protocol/*.proto
