package utils

import (
	"go_chain/block"
	"go_chain/transaction"
)

func NewGenesisBlock() *block.Block {
	block := &block.Block{}
	block.SetHash("Genesis")
	return block
}

func NewCoinbase(pubKey string, value uint64) *transaction.Transaction {
	tx := transaction.New()
	tx.AddInput(transaction.NewInput(0, "some hash"))
	tx.AddOutput(transaction.NewOutput("public key", 500))
	tx.CalcHash()
	return tx
}
