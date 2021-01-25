PROJECT_NAME = chain_proto
OUT_DIR = $(GOPATH)/src/$(PROJECT_NAME)/bin

# Go varibales
GOBIN = $(GOPATH)/bin
BUILD = GOOS=$(GOOS) go build

# protobuf variables
PROTO_DIR = gateway/proto
PB_OUT_DIR = gateway/gw
GOOGLEAPIS_DIR = gateway/proto/googleapis

ifeq ($(shell uname),Linux)
	export GOOS=linux 
endif

all: cli

cli: wallet server client

clean: 
	rm data/blockchain.db

wallet:
	$(BUILD) -o $(OUT_DIR)/wallet $(PROJECT_NAME)/cmd/wallet

server: 
	$(BUILD) -o $(OUT_DIR)/server $(PROJECT_NAME)/cmd/server

client: 
	$(BUILD) -o $(OUT_DIR)/client $(PROJECT_NAME)/cmd/client

proto:
	protoc -I $(PROTO_DIR) -I $(GOOGLEAPIS_DIR)\
	  --go_out $(PB_OUT_DIR) --go_opt paths=source_relative --plugin=$(GOBIN)/protoc-gen-go\
	  --go-grpc_out $(PB_OUT_DIR) --go-grpc_opt paths=source_relative --plugin=$(GOBIN)/protoc-gen-go-grpc\
	  --grpc-gateway_out $(PB_OUT_DIR) --plugin=$(GOBIN)/protoc-gen-grpc-gateway\
      --grpc-gateway_opt logtostderr=true \
      --grpc-gateway_opt paths=source_relative \
      --grpc-gateway_opt generate_unbound_methods=true \
	  $(PROTO_DIR)/*.proto

.PHONY: wallet server cli proto clean client
