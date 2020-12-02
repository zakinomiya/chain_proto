package miner

import (
	"go_chain/block"
	"go_chain/transaction"
	"go_chain/wallet"
	"sync"
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
	b := &MockBlockchain{[]*block.Block{}}
	w, _ := wallet.New()
	m := New(b, w)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go m.Start()
	time.AfterFunc(time.Second*5, func() {
		wg.Done()
	})

	wg.Wait()
}
