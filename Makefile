PROTO_PATH=/home/haxolotl/go/src/github.com/adithyabhatkajake/libe2c
GO_OUT_DIR=/home/haxolotl/go/src/

# Reference on how to use $^, $? and other automatic variables
# https://www.gnu.org/software/make/manual/html_node/Automatic-Variables.html

# Build all the protobuf files
proto: proto/config.proto \
	proto/crypto.proto \
	proto/e2c.proto \
	proto/net.proto \
	proto/e2c/blame.proto \
	proto/e2c/command.proto \
	proto/e2c/generic.proto \
	proto/e2c/proposal.proto
	@echo "Using Proto Path: ${PROTO_PATH}"
	@echo "Using Go Out Directory: ${GO_OUT_DIR}"
	# Build only changed protobuf definitions
	# Compiling Protobuf tips
	# https://jbrandhorst.com/post/go-protobuf-tips/
	protoc $? -I${PROTO_PATH} --go_out=:${GO_OUT_DIR}

gen-config: tools/genConfig.go
	go build -o tools/genConfig $^

gen-test-data: gen-config
	@echo "TODO: Generate a config for 10 nodes in testData directory"
	@echo "TODO: Generate a config for 3 nodes in testData directory"

rbc-replica: node/rbc/rbc_node.go
	go build -o node/rbc/rbc_replica $^

clean:
	@rm -rf tools/genConfig
	@rm -rf node/rbc/rbc_replica