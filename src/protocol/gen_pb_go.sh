# P.S PROTOBUF_DIR=~/protobuf/src
PROTOBUF_DIR=your protobuf dir
PROTOC_TOOL=${PROTOBUF_DIR}/protoc
${PROTOC_TOOL} --plugin=../../bin/protoc-gen-go --go_out=plugins=pepper,import_path=protocol:. --proto_path=${PROTOBUF_DIR}:. *.proto
