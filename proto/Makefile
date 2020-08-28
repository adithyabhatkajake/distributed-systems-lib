PROTO_PATH=/home/haxolotl/go/src/github.com/adithyabhatkajake/libchatter/proto
GO_OUT_DIR=/home/haxolotl/go/src/

# Reference on how to use $^, $? and other automatic variables
# https://www.gnu.org/software/make/manual/html_node/Automatic-Variables.html

# Build all the protobuf files
.PHONY: all e2c synchs apollo clean

all: e2c synchs apollo

e2c: crypto.proto \
	net.proto \
	e2c/config/node.proto \
	e2c/config/protocol.proto \
	e2c/config/client-config.proto \
	chain/command.proto \
	chain/block.proto \
	e2c/msg/proposal.proto \
	e2c/msg/blame.proto \
	e2c/msg/generic.proto \
	e2c/msg/clientAck.proto
	@echo "Using Proto Path: ${PROTO_PATH}"
	@echo "Using Go Out Directory: ${GO_OUT_DIR}"
	# Build only changed protobuf definitions
	# Compiling Protobuf tips
	# https://jbrandhorst.com/post/go-protobuf-tips/
	protoc $? -I${PROTO_PATH} --go_out=:${GO_OUT_DIR}

synchs: synchs/config/node.proto \
	synchs/config/protocol.proto \
	synchs/config/client-config.proto \
	synchs/msg/proposal.proto \
	synchs/msg/blame.proto \
	synchs/msg/generic.proto \
	synchs/msg/clientAck.proto \
	synchs/msg/vote.proto \
	synchs/msg/cert.proto
	@echo "Using Proto Path: ${PROTO_PATH}"
	@echo "Using Go Out Directory: ${GO_OUT_DIR}"
	# Build only changed protobuf definitions
	# Compiling Protobuf tips
	# https://jbrandhorst.com/post/go-protobuf-tips/
	protoc $? -I${PROTO_PATH} --go_out=:${GO_OUT_DIR}

apollo: apollo/config/node.proto \
	apollo/config/protocol.proto \
	apollo/config/client-config.proto \
	apollo/msg/proposal.proto \
	apollo/msg/blame.proto \
	apollo/msg/generic.proto \
	apollo/msg/clientAck.proto
	@echo "Using Proto Path: ${PROTO_PATH}"
	@echo "Using Go Out Directory: ${GO_OUT_DIR}"
	# Build only changed protobuf definitions
	# Compiling Protobuf tips
	# https://jbrandhorst.com/post/go-protobuf-tips/
	protoc $? -I${PROTO_PATH} --go_out=:${GO_OUT_DIR}

clean:
	@echo "Doing nothing for proto clean"