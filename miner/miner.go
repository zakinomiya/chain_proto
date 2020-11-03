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
	*wallet.Wallet
}

func New(bc blockchain.BlockchainInterface) *Miner {
	return &Miner{blockchain: bc}
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

func (m *Miner) CalcGenesis() *block.Block {
	gen := block.NewGenesisBlock()

	m.findNonce(gen, make(chan struct{}), 5)
	return gen
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
	runWorkers := func(workersCount uint32) {
		for i := uint32(0); i < workersCount; i++ {
			quit := make(chan struct{})

			go m.generateBlock(quit)
		}
	}
	runWorkers(uint32(runtime.NumCPU()))

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

func (m *Miner) findNonce(block *block.Block, quit chan struct{}, target uint32) bool {
	consecutiveZeros := strings.Repeat("0", int(target))

	for {
		select {
		case <-m.quit:
			return false
		default:
			//
		}
		hash := block.HashBlock()

		// TODO research on bitcoin compare algorithm is needed
		if strings.HasPrefix(fmt.Sprintf("%x", hash), consecutiveZeros) {
			block.Hash = hash
			log.Printf("info: Found a valid nonce: %v \n", block.Nonce)
			return true
		}
		block.IncrementNonce()
	}
}

func (m *Miner) generateBlock(quit chan struct{}) {
	log.Println("trace: Started generating a new block")

	for {
		select {
		case <-m.quit:
			log.Println("trace: Breaking OUTER: action=generateBlock")
			return
		default:
			//
		}

		coinbase := transaction.NewCoinbase(m.X.Bytes(), 25)
		block := m.blockchain.GenerateBlock(append([]*transaction.Transaction{coinbase}, m.transactionPool...))
		if m.findNonce(block, quit, block.Bits) {
			m.blockchain.AddBlock(block)
			log.Printf("info: %x \n", block.Hash)
			log.Println("debug: Closing the quit channel")
			m.Restart()
		}
	}
}
