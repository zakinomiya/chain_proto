package utils

import (
	"go_chain/block"
	"go_chain/transaction"
	"math/rand"
)

func NewGenesisBlock() *block.Block {
	block := block.New()
	genesis := NewCoinbase([]byte("This is Minimum Viable Blockchain"), 25)
	block.SetTranscations([]*transaction.Transaction{genesis})
	block.CalcMerkleTree()
	return block
}

func NewCoinbase(pubKey []byte, value uint64) *transaction.Transaction {
	tx := transaction.New()
	tx.AddOutput(transaction.NewOutput())
	tx.CalcHash()
	return tx
}

/// Pseudo random uint32
func RandomUint32() uint32 {
	return rand.Uint32()
}
