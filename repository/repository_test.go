package repository

import (
	"go_chain/block"
	"go_chain/transaction"
	"testing"
	"time"
)

func TestFromBlock(t *testing.T) {
	// r := New("./data/blockchain.db", "sqlite3")
	b := &block.Block{}
	b.Hash = [32]byte{}
	b.Height = 100
	b.MerkleRoot = []byte("Merkle Root")
	b.PrevBlockHash = [32]byte{}
	b.Timestamp = uint32(time.Now().Unix())
	b.Bits = 5
	b.ExtraNonce = 123456
	b.Nonce = 123456789
	b.Transactions = []*transaction.Transaction{}

	// blockModel := &BlockModel{}
	// r.fromBlock(b, blockModel)

	// t.Errorf("%+v", blockModel)
}
