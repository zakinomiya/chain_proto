package miner

import (
	"fmt"
	"go_chain/block"
	"go_chain/blockchain"
	"go_chain/transaction"
	"go_chain/wallet"
	"log"
	"strings"
	"sync"
	"time"
)

type State int

const (
	unknown State = iota
	running
	stopping
	stopped
	restarting
)

type event int

const (
	none event = iota + 1
	found
	interrupted
)

const (
	defaultMaxWorkersNum = 5
)

func (s State) String() string {
	switch s {
	case unknown:
		return "UNKNOWN"
	case stopping:
		return "STOPPING"
	case stopped:
		return "STOPPED"
	case restarting:
		return "RESTARTING"
	default:
		return "UNKNOWN"
	}
}

type Blockchain interface {
	CurrentBlockHeight() uint32
	Difficulty() uint32
	LatestBlock() *block.Block
	AddBlock(block *block.Block) bool
	Subscribe(key string) <-chan blockchain.BlockchainEvents
	Unsubscribe(key string)
}

type Miner struct {
	enabled          bool
	concurrent       bool
	maxWorkersNum    int
	wg               *sync.WaitGroup
	blockchainEvents <-chan blockchain.BlockchainEvents
	exit             chan struct{}
	transactionPool  []*transaction.Transaction
	blockchain       Blockchain
	minerWallet      *wallet.Wallet
}

func New(bc Blockchain, w *wallet.Wallet, enabled bool, concurrent bool, maxWorkersNum int) *Miner {
	// listening blockchain updates
	ch := bc.Subscribe("Miner")

	return &Miner{
		enabled:          enabled,
		concurrent:       concurrent,
		maxWorkersNum:    maxWorkersNum,
		wg:               &sync.WaitGroup{},
		blockchainEvents: ch,
		exit:             make(chan struct{}),
		transactionPool:  []*transaction.Transaction{},
		blockchain:       bc,
		minerWallet:      w,
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
	events := []chan blockchain.BlockchainEvents{}

	work := func() {
		for i := 0; i < workersNum; i++ {
			q := make(chan struct{})
			e := make(chan blockchain.BlockchainEvents)
			workers = append(workers, q)
			events = append(events, e)
			go m.worker(q, e)
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
		case event := <-m.blockchainEvents:
			log.Println("debug: received a new event")
			for _, e := range events {
				e <- event
			}
		}
	}
}

func (m *Miner) worker(quit chan struct{}, eventStream <-chan blockchain.BlockchainEvents) {
	block := m.generateBlock()
	consecutiveZeros := strings.Repeat("0", int(block.Bits))

	for {
		select {
		case <-quit:
			log.Println("debug: received signal. quit working")
			m.wg.Done()
			return
		case event := <-eventStream:
			log.Printf("debug: new event received: event=%s\n", event)
			if event == blockchain.NewBlock {
				log.Println("debug: update block")
				block = m.generateBlock()
			}
		default:
			//
		}

		if hash := block.HashBlock(); strings.HasPrefix(fmt.Sprintf("%x", hash), consecutiveZeros) {
			block.Hash = hash
			log.Printf("info: Found a valid nonce: %v \n", block.Nonce)
			m.blockchain.AddBlock(block)
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
