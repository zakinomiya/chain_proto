PROJECT_NAME = go_chain
OUT_DIR = bin

.PHONY: wallet

wallet:
	go build -o $(OUT_DIR)/wallet $(PROJECT_NAME)/cmd 