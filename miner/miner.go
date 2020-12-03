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

func (m *Miner) Start() error {
	log.Println("info: Starting mining process")

	done := make(chan struct{}, 0)
	m.done = done
	m.state = running

	go m.mining()

	for {
		select {
		case <-done:
			if m.state == running {
				m.state = restarting
				m.Restart()
				return nil
			} else if m.state == stopping {
				m.state = stopped
				time.Sleep(time.Second * 1)
				log.Println("info: Mining server gracefully stopped")
				return nil
			}

			return nil
		default:
			//
		}
	}
}

func (m *Miner) Stop() {
	log.Println("info: Stopping mining process")
	m.state = stopping
	m.interrupt()
}

func (m *Miner) ServiceName() string {
	return "Miner"
}

func (m *Miner) Restart() {
	log.Println("info: Restarting Miner")
	time.Sleep(time.Second * miningWaitSecond)
	m.Start()
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

func (m *Miner) mining() {
	log.Println("info: Started mining")
	found := make(chan struct{}, 0)
	for i := 0; i < maxMiningNum; i++ {
		b := m.generateBlock()
		go m.findNonce(found, b)
	}

	for {
		select {
		case <-found:
			log.Println("info: Someone found a nonce")
			m.interrupt()
			return
		default:
			//
		}
	}
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
			close(found)
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
