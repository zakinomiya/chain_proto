package utils

import (
	"go_chain/block"
	"go_chain/transaction"
	"math/rand"
)

func NewGenesisBlock() *block.Block {
	b := block.New()
	h := block.NewHeader()
	h.Bits = 5
	h.Nonce = 129
	h.MerkleRoot = []byte("merkle")
	h.PrevBlockHash = [32]byte{}
	// coinbase := NewCoinbase([]byte("This is Minimum Viable Blockchain"), 25)
	b.BlockHeader = h
	b.Height = 1000
	// b.Transactions = []*transaction.Transaction{coinbase}
	b.SetExtranNonce()
	return b
}

func NewCoinbase(pubKey []byte, value uint64) *transaction.Transaction {
	tx := transaction.New()
	tx.CalcHash()
	return tx
}

/// Pseudo random uint32
func RandomUint32() uint32 {
	return rand.Uint32()
}
