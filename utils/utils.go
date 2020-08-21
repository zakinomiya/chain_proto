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
	// tx := transaction.New()
	// tx.AddOutput(transaction.NewOutput(0, pubKey, value))
	// tx.CalcHash()
	// tx.Sign()
	// return tx
}
