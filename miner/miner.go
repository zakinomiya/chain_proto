package miner

import (
	"chain_proto/block"
	"chain_proto/blockchain"
	"chain_proto/transaction"
	"chain_proto/wallet"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

const (
	defaultMaxWorkersNum = 5
	miningWaitTime       = time.Second * 10
)

type Blockchain interface {
	CurrentBlockHeight() uint32
	Difficulty() uint32
	LatestBlock() *block.Block
	AddBlock(block *block.Block) bool
	Subscribe(key string) <-chan blockchain.BlockchainEvents
	Unsubscribe(key string)
}

type Miner struct {
	enabled         bool
	concurrent      bool
	maxWorkersNum   int
	wg              *sync.WaitGroup
	exit            chan struct{}
	transactionPool []*transaction.Transaction
	blockchain      Blockchain
	minerWallet     *wallet.Wallet
}

func New(bc Blockchain, w *wallet.Wallet, enabled bool, concurrent bool, maxWorkersNum int) *Miner {
	return &Miner{
		enabled:         enabled,
		concurrent:      concurrent,
		maxWorkersNum:   maxWorkersNum,
		wg:              &sync.WaitGroup{},
		exit:            make(chan struct{}),
		transactionPool: []*transaction.Transaction{},
		blockchain:      bc,
		minerWallet:     w,
	}
}

func (m *Miner) Start() error {
	log.Println("info: Starting mining process")

	workersNum := defaultMaxWorkersNum

	if !m.enabled {
		log.Println("info: miner")
		return nil
	}

	if !m.concurrent {
		log.Println("info: [Miner] running in single mode")
		workersNum = 1
	} else {
		log.Println("info: [Miner] running in conccurent mode")
		workersNum = m.maxWorkersNum
	}

	go m.mining(workersNum)
	return nil
}

func (m *Miner) Stop() {
	log.Println("info: Stopping mining process")
	m.interrupt()
	time.Sleep(1 * time.Second)
}

func (m *Miner) ServiceName() string {
	return "Miner"
}

func (m *Miner) interrupt() {
	m.exit <- struct{}{}
}

func (m *Miner) AddTransaction(tx *transaction.Transaction) {
	m.transactionPool = append(m.transactionPool, tx)
}

func (m *Miner) mining(workersNum int) {
	log.Println("info: Started mining")
	workers := []chan struct{}{}

	work := func() {
		for i := 0; i < workersNum; i++ {
			q := make(chan struct{})
			workers = append(workers, q)
			go m.worker(q)
			m.wg.Add(1)
		}
	}

	go work()

	for {
		select {
		case <-m.exit:
			log.Println("info: received exit signal")
			for _, q := range workers {
				close(q)
			}

			m.wg.Wait()
			return
		}
	}
}

func (m *Miner) worker(quit chan struct{}) {
	var block *block.Block
	var consecutiveZeros string

	for {
		if block == nil {
			block = m.generateBlock()
			consecutiveZeros = strings.Repeat("0", int(block.Bits))
		}

		// Check blockchain update every 10000 calculations.
		if block.Nonce%10000 == 0 && m.blockchain.CurrentBlockHeight()+1 != block.Height {
			log.Printf("debug: new block already added. updating target block")
			block = nil
			time.Sleep(miningWaitTime)
			continue
		}

		select {
		case <-quit:
			log.Println("debug: received signal. quit working")
			m.wg.Done()
			return
		default:
			//
		}

		if hash := block.HashBlock(); strings.HasPrefix(fmt.Sprintf("%x", hash), consecutiveZeros) {
			block.Hash = hash
			log.Printf("info: Found a valid nonce: %v \n", block.Nonce)
			m.blockchain.AddBlock(block)
			block = nil
			time.Sleep(miningWaitTime)
			continue
		}
		block.IncrementNonce()
	}
}

func (m *Miner) generateBlock() *block.Block {
	coinbase := transaction.NewCoinbase(m.minerWallet, 25)
	txs := append([]*transaction.Transaction{coinbase}, m.GetPooledTransactions(10)...)
	block := block.New(m.blockchain.CurrentBlockHeight()+1, m.blockchain.Difficulty(), m.blockchain.LatestBlock().Hash, txs)
	return block
}
