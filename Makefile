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