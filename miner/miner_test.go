package miner

import (
	"go_chain/block"
	"go_chain/common"
	"go_chain/transaction"
	"testing"
	"time"
)

func TestMining(t *testing.T) {
	m := &Miner{}
	c := transaction.NewCoinbase([]byte("hello world"), 25)
	h := block.NewHeader()
	h.PrevBlockHash = [32]byte{}
	h.Timestamp = uint32(time.Now().Unix())
	h.Bits = 5
	h.Nonce = 0

	block := &block.Block{
		Height:       0,
		Hash:         [32]byte{},
		Transactions: []*transaction.Transaction{c},
		ExtraNonce:   common.RandomUint32(),
		BlockHeader:  h,
	}
	block.CalcMerkleTree()
	if m.findNonce(block, make(chan struct{}, 0), 5) {
		t.Errorf("height=%d, hash=%x, transaction=%#v, extraNonce=%d, prevBlockHash=%x, merkelRoot=%x, timestamp=%d, bits=%d, nonce=%d", block.Height, block.Hash, block.Transactions[0], block.ExtraNonce, block.PrevBlockHash, block.MerkleRoot, block.Timestamp, block.Bits, block.Nonce)
	}
}
