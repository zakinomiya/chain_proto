package utils

import (
	"go_chain/block"
	"go_chain/transaction"
	"math/rand"
)

func NewGenesisBlock() *block.Block {
	block := &block.Block{}
	return block
}

func NewCoinbase(pubKey []byte, value uint64) *transaction.Transaction {
	tx := transaction.New()
	tx.AddInput(transaction.NewInput(0, "some hash"))
	tx.AddOutput(transaction.NewOutput("public key", 500))
	tx.CalcHash()
	return tx
}

/// Pseudo random uint32
func RandomUint32() uint32 {
	return rand.Uint32()
}
