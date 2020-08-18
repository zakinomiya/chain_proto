package utils

import "go_chain/block"

func NewGenesisBlock() *block.Block {
	block := &block.Block{}
	block.SetHash("Genesis")
	block.SetAmount(100)
	return block
}
