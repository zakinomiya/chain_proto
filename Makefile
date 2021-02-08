PROJECT_NAME = chain_proto
OUT_DIR = $(GOPATH)/src/$(PROJECT_NAME)/bin

# Go varibales
GOBIN = $(GOPATH)/bin
FLAGS = 
BUILD = GOOS=$(GOOS) go build $(FLAGS)

# Docker variables
IMAGE_NAME = "chain"

# protobuf variables
PROTO_DIR = gateway/proto
PB_OUT_DIR = gateway/gw
GOOGLEAPIS_DIR = gateway/proto/googleapis

ifeq ($(shell uname),Linux)
	export GOOS=linux 
	export FLAGS=-ldflags="-extldflags=-static" -tags sqlite_omit_load_extension
endif

ifeq ($(shell uname),Darwin)
	export GOOS=darwin 
endif

all: cli

cli: wallet server client

clean: 
	rm data/blockchain.db

image:
	docker build -t "${IMAGE_NAME}" .

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
	  --govalidators_out=gateway/gw --plugin=$(GOBIN)/protoc-gen-govalidators\
      --grpc-gateway_opt logtostderr=true \
      --grpc-gateway_opt paths=source_relative \
      --grpc-gateway_opt generate_unbound_methods=true \
	  $(PROTO_DIR)/*.proto

.PHONY: wallet server cli proto clean client
