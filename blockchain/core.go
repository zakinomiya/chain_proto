package blockchain

import (
	"chain_proto/block"
	"chain_proto/db/repository"
	"chain_proto/transaction"
	"fmt"
	"log"
	"os"
	"sync"
)

type Miner interface {
	AddTransaction(tx *transaction.Transaction) error
}

// Blockchain is a struct of the chain
type Blockchain struct {
	mutex      sync.Mutex
	chainID    uint32
	height     uint32
	blocks     []*block.Block
	repository *repository.Repository
	miner      Miner
}

var blockchain *Blockchain
var once sync.Once

// New returns a new blockchain
func New(repository *repository.Repository) *Blockchain {
	blockchain = &Blockchain{
		repository: repository,
	}
	return blockchain
}

func initializeBlockchain() error {
	b, err := blockchain.repository.Block.GetLatest()

	if err == repository.ErrNotFound {
		log.Println("info: No blocks found in the db. Creating the genesis block")
		genesis, err := block.NewGenesisBlock()
		if err != nil {
			return err
		}

		if err := blockchain.AddBlock(genesis); err != nil {
			return fmt.Errorf("error: failed to add the genesis block. err=%s\n", err)
		}
		return nil
	} else {
		return err
	}

	log.Println("info: Block record found in the db. Restoring the blockchain")
	blockchain.height = b.Height
	blockchain.blocks = []*block.Block{b}
	return nil
}

// ServiceName returns its service name
func (bc *Blockchain) ServiceName() string {
	return "Blockchain"
}

// Start intialises and starts the blockchain
// There must be only one blockchain during the runtime.
func (bc *Blockchain) Start() error {
	once.Do(func() {
		if err := initializeBlockchain(); err != nil {
			log.Println("error:", err)
			log.Println("Exiting the process...")
			os.Exit(1)
		}
	})
	return nil
}

// Stop stops the blockchain gracefully.
// TODO inplement Stop
func (bc *Blockchain) Stop() {
	return
}

// Difficulty returns the next mining target.
// TODO implement Difficulty
func (bc *Blockchain) Difficulty() uint32 {
	return 5
}
