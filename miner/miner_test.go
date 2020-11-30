package miner

import (
	"encoding/json"
	"go_chain/block"
	"go_chain/common"
	"go_chain/transaction"
	"go_chain/wallet"
	"io/ioutil"
	"testing"
	"time"
)

type MockBlockchain struct {
	blocks []*block.Block
}

func newMock() *MockBlockchain {
	return &MockBlockchain{blocks: []*block.Block{block.New(1, 5, [32]byte{}, make([]*transaction.Transaction, 0))}}
}

func (bc *MockBlockchain) CurrentBlockHeight() uint32 {
	return uint32(len(bc.blocks))
}

func (bc *MockBlockchain) Difficulty() uint32 {
	return 5
}

func (bc *MockBlockchain) LatestBlock() *block.Block {
	return block.New(bc.CurrentBlockHeight()+1, bc.Difficulty(), [32]byte{}, make([]*transaction.Transaction, 0))
}

func (bc *MockBlockchain) AddBlock(block *block.Block) bool {
	bc.blocks = append(bc.blocks, block)
	return true
}

func (bc *MockBlockchain) GetPooledTransactions(num int) []*transaction.Transaction {
	txs := []*transaction.Transaction{}
	for i := 0; i < num; i++ {
		txs = append(txs, transaction.New())
	}

	return txs
}

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
	if m.findNonce(block) {
		j, _ := json.Marshal(block)
		t.Log(string(j))
		ioutil.WriteFile("block.json", j, 0400)
	}
}

func TestConcurrentMining(t *testing.T) {
	b := &MockBlockchain{[]*block.Block{}}
	w, _ := wallet.New()

	m := New(b, w)
	if err := m.Start(); err != nil {
		t.Logf("error: %x", err)
	}
	time.Sleep(time.Second * 10)
}
