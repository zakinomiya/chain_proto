package miner

import (
	"fmt"
	"go_chain/block"
	"go_chain/blockchain"
	"go_chain/transaction"
	"go_chain/utils"
	"log"
	"runtime"
	"sync"
	"time"
)

type Miner struct {
	transactionPool []*transaction.Transaction
	quit            chan struct{}
	blockLock       sync.Mutex
	wg              sync.WaitGroup
	workers         []chan struct{}
	pubKey          []byte /// TODO read from config
	blockchain      blockchain.BlockchainInterface
}

func New(bc blockchain.BlockchainInterface) *Miner {
	return &Miner{blockchain: bc}
}

func (m *Miner) Start() error {
	m.quit = make(chan struct{})
	m.wg.Add(1)
	go m.mining()

	log.Println("Mining process started")
	return nil
}

func (m *Miner) Stop() {
	close(m.quit)
	m.wg.Wait()

	log.Println("Mining process stopped")
}

func (m *Miner) ServiceName() string {
	return "Miner"
}

func (m *Miner) AddTransaction(tx *transaction.Transaction) {
	m.transactionPool = append(m.transactionPool, tx)
}

func (m *Miner) mining() {
	runWorkers := func(workersCount uint32) {
		fmt.Printf("CPU NUM: %v \n", workersCount)
		for i := uint32(0); i < workersCount; i++ {
			fmt.Printf("worker no. %v \n", i)
			quit := make(chan struct{})
			m.workers = append(m.workers)
			m.wg.Add(1)

			go m.generateBlock(quit)
		}
	}
	runWorkers(uint32(runtime.NumCPU()))

	select {
	case <-m.quit:
		for _, w := range m.workers {
			close(w)
		}
	}
}

func (m *Miner) findNonce(block *block.Block, quit chan struct{}) bool {

	time.Sleep(time.Second * 1)

	if utils.RandomUint32() > 50 {
		quit <- struct{}{}
	}

	return true
}

func (m *Miner) generateBlock(quit chan struct{}) {
	log.Println("Started generating a new block")

OUTER:
	for {

		select {
		case <-quit:
			break OUTER
		default:
			//
		}

		block := block.New()
		coinbase := utils.NewCoinbase(m.pubKey, 25)
		block.SetTranscations(append([]*transaction.Transaction{coinbase}, m.transactionPool...))
		m.findNonce(block, quit)
	}
}
