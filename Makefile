PROJECT_NAME = go_chain
OUT_DIR = bin

.PHONY: wallet miner cli

all: cli

cli: wallet miner

wallet:
	go build -o $(OUT_DIR)/wallet $(PROJECT_NAME)/cmd/wallet

miner: 
	go build -o $(OUT_DIR)/miner $(PROJECT_NAME)/cmd/miner