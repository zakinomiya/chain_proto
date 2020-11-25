package miner

import (
	"encoding/json"
	"go_chain/block"
	"go_chain/common"
	"go_chain/transaction"
	"go_chain/wallet"
	"io/ioutil"
	"testing"
)

func TestMining(t *testing.T) {
	w, _ := wallet.New()
	m := &Miner{minerWallet: w}
	c := transaction.NewCoinbase(w, 25)
	h := block.NewHeader()
	h.PrevBlockHash = [32]byte{}
	h.Timestamp = common.Timestamp()
	h.Bits = 5
	h.Nonce = 0

	block := &block.Block{
		Height:       0,
		Hash:         [32]byte{},
		Transactions: []*transaction.Transaction{c},
		ExtraNonce:   common.RandomUint32(),
		BlockHeader:  h,
	}
	block.SetMerkleRoot()
	if m.findNonce(block, make(chan struct{}, 0)) {
		j, _ := json.Marshal(block)
		t.Log(string(j))
		ioutil.WriteFile("block.json", j, 0400)
	}
}
