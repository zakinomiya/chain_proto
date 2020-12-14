PROJECT_NAME = go_chain
OUT_DIR = bin

.PHONY: wallet server cli

all: cli

cli: wallet server

clean: 
	rm data/blockchain.db

wallet:
	go build -o $(OUT_DIR)/wallet $(PROJECT_NAME)/cmd/wallet

server: 
	go build -o $(OUT_DIR)/server $(PROJECT_NAME)/cmd/server

# TODO define cmd to compile proto files 
# here is a command I tried
# protoc -I . -Igoogleapis \
#    --go_out ./ --go_opt paths=source_relative --plugin=$GOPATH/bin/protoc-gen-go  \
#    --go-grpc_out ./ --go-grpc_opt paths=source_relative --plugin=$GOPATH/bin/protoc-gen-grpc-gateway --plugin=$GOPATH/bin/protoc-gen-go-grpc     \
#    http_service.proto