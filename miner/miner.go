package miner

import (
	"context"
	"fmt"
	"go_chain/block"
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

const (
	maxMiningNum     = 10
	miningWaitSecond = 3
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
}

type Miner struct {
	mux             *sync.Mutex
	state           State
	done            chan struct{}
	miningCtx       context.Context
	transactionPool []*transaction.Transaction
	blockchain      Blockchain
	minerWallet     *wallet.Wallet
}

func New(bc Blockchain, w *wallet.Wallet) *Miner {
	return &Miner{mux: &sync.Mutex{}, blockchain: bc, minerWallet: w}
}

func (m *Miner) Start() error {
	log.Println("info: Starting mining process")

	done := make(chan struct{}, 0)
	m.done = done

	m.state = running

	if err := m.mining(done); err != nil {
		return err
	}
	log.Println("info: Stopped mining process")
	return nil
}

func (m *Miner) Stop() {
	log.Println("info: Stopping mining process")
	m.mux.Lock()
	defer m.mux.Unlock()
	m.state = stopping
	m.interrupt()
	time.Sleep(1 * time.Second)
}

func (m *Miner) ServiceName() string {
	return "Miner"
}

func (m *Miner) Restart() {
	m.state = restarting
	log.Println("info: Restarting Miner")
	time.Sleep(time.Second * miningWaitSecond)
	go m.Start()
}

func (m *Miner) Status() string {
	return m.state.String()
}

func (m *Miner) interrupt() {
	m.done <- struct{}{}
}

func (m *Miner) AddTransaction(tx *transaction.Transaction) {
	m.transactionPool = append(m.transactionPool, tx)
}

// TODO wait until a new block is stored in the chain
func (m *Miner) mining(done chan struct{}) error {
	log.Println("info: Started mining")
	found := make(chan struct{}, 0)
	m.startWorkers(found)

	for {
		select {
		case <-done:
			log.Println("debug: mining interrupted")
			m.stopWorkers(found)
			return nil
		case <-found:
			log.Println("info: Someone found a nonce")
			m.startWorkers(found)
			continue
		default:
			//
		}
	}
}

func (m *Miner) startWorkers(found chan struct{}) {
	m.mux.Lock()
	defer m.mux.Unlock()
	for i := 0; i < maxMiningNum; i++ {
		b := m.generateBlock()
		go m.findNonce(found, b)
	}
	m.state = running
}

func (m *Miner) stopWorkers(found chan struct{}) {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.state = stopping
	close(found)
	time.Sleep(time.Second * miningWaitSecond)
	m.state = stopped
}

func (m *Miner) findNonce(found chan struct{}, block *block.Block) {
	consecutiveZeros := strings.Repeat("0", int(block.Bits))

	for {
		select {
		case <-found:
			return
		default:
			//
		}

		hash := block.HashBlock()

		if strings.HasPrefix(fmt.Sprintf("%x", hash), consecutiveZeros) {
			block.Hash = hash
			log.Printf("info: Found a valid nonce: %v \n", block.Nonce)
			found <- struct{}{}
			m.stopWorkers(found)
			m.blockchain.AddBlock(block)
			return
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
