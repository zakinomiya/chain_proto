package miner

import (
	"encoding/json"
	"go_chain/block"
	"go_chain/transaction"
	"go_chain/wallet"
	"log"
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
	j, _ := json.Marshal(block)
	log.Printf("info: block=%s", string(j))

	bc.blocks = append(bc.blocks, block)
	log.Printf("info: now blockchain height is %d\n", len(bc.blocks))
	return true
}

func TestMining(t *testing.T) {
	b := &MockBlockchain{[]*block.Block{}}
	w, _ := wallet.RestoreWallet("58898c79caf4a77a4aa10b4b9bad7d07f7e7c1842204be352a65d87f71277137")
	m := New(b, w)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go m.Start()
	time.AfterFunc(time.Second*5, func() {
		wg.Done()
	})

	wg.Wait()
}
