package miner

import (
	"fmt"
	"go_chain/block"
	"go_chain/blockchain"
	"go_chain/transaction"
	"go_chain/wallet"
	"log"
	"runtime"
	"strings"
)

type Miner struct {
	quit            chan struct{}
	transactionPool []*transaction.Transaction
	blockchain      blockchain.BlockchainInterface
	minerWallet     *wallet.Wallet
}

func New(bc blockchain.BlockchainInterface, w *wallet.Wallet) *Miner {
	return &Miner{blockchain: bc, minerWallet: w}
}

func (m *Miner) Start() error {
	if m.quit != nil {
		log.Println("debug: Mining process is already running")
		return nil
	}

	m.quit = make(chan struct{})
	go m.mining()

	log.Println("debug: Mining process started")
	return nil
}

func (m *Miner) Stop() {
	if m.quit == nil {
		log.Println("debug: No mining process is running")
		return
	}
	close(m.quit)
	m.quit = nil
	log.Println("debug: Mining process stopped")
}

func (m *Miner) ServiceName() string {
	return "Miner"
}

func (m *Miner) Restart() {
	log.Println("debug: Restarting Miner")
	m.Stop()
	m.Start()
}

func (m *Miner) AddTransaction(tx *transaction.Transaction) {
	m.transactionPool = append(m.transactionPool, tx)
}

func (m *Miner) mining() {
	runWorkers := func(workersCount int) {
		for i := 0; i < workersCount; i++ {
			b := m.generateBlock()
			go m.findNonce(b)
		}
	}
	runWorkers(runtime.NumCPU())

	for {
		select {
		case <-m.quit:
			log.Println("debug: Someone closed the quit channel")
			return
		default:
			//
		}
	}
}

func (m *Miner) findNonce(block *block.Block) bool {
	log.Println("Started mining")
	consecutiveZeros := strings.Repeat("0", int(block.Bits))

	for {
		select {
		case <-m.quit:
			return false
		default:
			//
		}
		hash := block.HashBlock()

		if strings.HasPrefix(fmt.Sprintf("%x", hash), consecutiveZeros) {
			block.Hash = hash
			log.Printf("info: Found a valid nonce: %v \n", block.Nonce)
			return true
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
