package blockchain

import (
	"chain_proto/block"
	"chain_proto/db/repository"
	"chain_proto/gateway"
	"chain_proto/transaction"
	"log"
	"os"
	"sync"
)

type Miner interface {
	AddTransaction(tx *transaction.Transaction)
}

// Blockchain is a struct of the chain
type Blockchain struct {
	chainID    string
	blocks     []*block.Block
	repository *repository.Repository
	miner      Miner
	client     *gateway.Client
}

var blockchain *Blockchain
var once sync.Once

// New returns a new blockchain
func New(repository *repository.Repository) *Blockchain {
	blockchain = &Blockchain{
		repository: repository,
		blocks:     make([]*block.Block, 0),
	}
	return blockchain
}

func (bc *Blockchain) SetMiner(m Miner) {
	bc.miner = m
}

func initializeBlockchain() error {
	b, err := blockchain.repository.Block.GetLatest()
	if err != nil {
		if err != repository.ErrNotFound {
			return err
		}

		log.Println("info: No blocks found in the db. Creating the genesis block")
		return blockchain.genesis()
	}

	blockchain.blocks = append(blockchain.blocks, b)
	return nil
}

func (bc *Blockchain) genesis() error {
	gen, err := block.NewGenesisBlock()
	if err != nil {
		return err
	}
	return blockchain.AddBlock(gen)
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
			log.Println("error: failed to initialise the blockchain", err)
			log.Println("Exiting the process...")
			os.Exit(1)
		}
		log.Println("info: finished initialising the blockchain")
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
