package miner

import (
	"context"
	"fmt"
	"go_chain/block"
	"go_chain/blockchain"
	"go_chain/transaction"
	"go_chain/wallet"
	"log"
	"strings"
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

func (s *State) Status(status string) int {
	switch status {
	case "Running":
		return 1
	case "Stopping":
		return 2
	case "Stopped":
		return 3
	default:
		return 0
	}
}

type Miner struct {
	state           State
	done            chan struct{}
	miningCtx       context.Context
	transactionPool []*transaction.Transaction
	blockchain      blockchain.BlockchainInterface
	minerWallet     *wallet.Wallet
}

func New(bc blockchain.BlockchainInterface, w *wallet.Wallet) *Miner {
	return &Miner{blockchain: bc, minerWallet: w}
}

func (m *Miner) Start() {
	log.Println("info: Starting mining process")

	done := make(chan struct{}, 0)
	m.done = done
	m.state = running

	go m.mining(done)

	for {
		select {
		case <-done:
			log.Println("Mining process stopped")
			if m.state == running {
				m.state = restarting
				m.Restart()
				return
			} else if m.state == stopping {
				m.state = stopped
				log.Println("Mining server gracefully stopped")
				return
			}
		default:
			//
		}
	}
}

func (m *Miner) Stop() {
	log.Println("info: Stopping mining procss")
	m.state = stopping
	m.done <- struct{}{}
	time.Sleep(time.Second * 2)
}

func (m *Miner) ServiceName() string {
	return "Miner"
}

func (m *Miner) Restart() {
	log.Println("debug: Restarting Miner")
	m.Start()
}

func (m *Miner) AddTransaction(tx *transaction.Transaction) {
	m.transactionPool = append(m.transactionPool, tx)
}

func (m *Miner) mining(done chan struct{}) {
	found := make(chan struct{}, 1)
	for i := 0; i < 10; i++ {
		b := m.generateBlock()
		go m.findNonce(found, b)
	}

	for {
		select {
		case <-found:
			log.Println("debug: Someone found a nonce")
			done <- struct{}{}
			return
		default:
			//
		}
	}
}

func (m *Miner) findNonce(found chan struct{}, block *block.Block) {
	log.Println("Started mining")
	consecutiveZeros := strings.Repeat("0", int(block.Bits))

	for {
		select {
		case <-found:
			log.Println("info: exiting findNonce")
			return
		default:
			//
		}

		hash := block.HashBlock()

		if strings.HasPrefix(fmt.Sprintf("%x", hash), consecutiveZeros) {
			block.Hash = hash
			log.Printf("info: Found a valid nonce: %v \n", block.Nonce)
			found <- struct{}{}
			close(found)
			return
		}
		block.IncrementNonce()
	}
}

func (m *Miner) generateBlock() *block.Block {
	log.Println("trace: Started generating a new block")
	coinbase := transaction.NewCoinbase(m.minerWallet, 25)
	txs := append([]*transaction.Transaction{coinbase}, m.blockchain.GetPooledTransactions(10)...)
	block := block.New(m.blockchain.CurrentBlockHeight(), m.blockchain.Difficulty(), m.blockchain.LatestBlock().Hash, txs)
	return block
}
