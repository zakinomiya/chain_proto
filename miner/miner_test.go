package miner

import (
	"encoding/json"
	"chain_proto/block"
	"chain_proto/blockchain"
	"chain_proto/config"
	"chain_proto/transaction"
	"chain_proto/wallet"
	"log"
	"sync"
	"testing"
	"time"
)

type MockBlockchain struct {
	blocks        []*block.Block
	subscriptions map[string]chan blockchain.BlockchainEvents
}

func newMock() *MockBlockchain {
	return &MockBlockchain{
		blocks:        []*block.Block{block.New(1, 5, [32]byte{}, make([]*transaction.Transaction, 0))},
		subscriptions: make(map[string]chan blockchain.BlockchainEvents),
	}
}

func (bc *MockBlockchain) CurrentBlockHeight() uint32 {
	return uint32(len(bc.blocks))
}

func (bc *MockBlockchain) Difficulty() uint32 {
	return 5
}

func (bc *MockBlockchain) LatestBlock() *block.Block {
	return bc.blocks[len(bc.blocks)-1]
}

func (bc *MockBlockchain) AddBlock(block *block.Block) bool {
	j, _ := json.Marshal(block)
	log.Printf("info: block=%s", string(j))

	bc.blocks = append(bc.blocks, block)
	log.Printf("info: now blockchain height is %d\n", block.Height)

	bc.SendEvent(blockchain.NewBlock)
	return true
}

func (bc *MockBlockchain) Subscribe(key string) <-chan blockchain.BlockchainEvents {
	ch := make(chan blockchain.BlockchainEvents)
	bc.subscriptions[key] = ch
	return ch
}

func (bc *MockBlockchain) Unsubscribe(key string) {
	delete(bc.subscriptions, key)
}

func (bc *MockBlockchain) SendEvent(eventName blockchain.BlockchainEvents) {
	for key, ch := range bc.subscriptions {
		log.Printf("debug: sending event(%s) to the subsctiption(%s)\n", eventName, key)
		ch <- eventName
	}
}

func TestMining(t *testing.T) {
	b := newMock()
	w, _ := wallet.RestoreWallet("58898c79caf4a77a4aa10b4b9bad7d07f7e7c1842204be352a65d87f71277137")
	m := New(b, w, config.Config.Enabled, config.Config.Concurrent, config.Config.MaxWorkersNum)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go m.Start()
	time.AfterFunc(time.Second*5, func() {
		m.Stop()
		wg.Done()
	})

	wg.Wait()
}
