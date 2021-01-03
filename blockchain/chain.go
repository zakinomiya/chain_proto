package blockchain

import (
	"chain_proto/block"
	"chain_proto/db/repository"
	"fmt"
	"log"
	"os"
	"sync"
)

// Blockchain is a struct of the chain
type Blockchain struct {
	mutex      sync.Mutex
	chainID    uint32
	height     uint32
	blocks     []*block.Block
	repository *repository.Repository
}

var blockchain *Blockchain
var once sync.Once

// New returns a new blockchain
func New(chainID uint32, repository *repository.Repository) *Blockchain {
	blockchain = &Blockchain{
		chainID:    chainID,
		repository: repository,
	}
	return blockchain
}

func initializeBlockchain() error {
	b, err := blockchain.repository.Block.GetLatestBlock()
	if err != nil {
		log.Println("error: Failed to initialise blockchain")
		return err
	}

	if b == nil {
		log.Println("info: No blocks found in the db. Creating the genesis block")
		genesis, err := block.NewGenesisBlock()
		if err != nil {
			return err
		}

		if !blockchain.AddBlock(genesis) {
			return fmt.Errorf("error: failed to add the genesis block")
		}
		return nil
	}

	log.Println("info: Block record found in the db. Restoring the blockchain")
	blockchain.height = b.Height
	blockchain.ReplaceBlocks([]*block.Block{b})
	return nil
}

// ServiceName returns its service name
func (bc *Blockchain) ServiceName() string {
	return "Blockchain"
}

// Start starts intialises and starts the blockchain
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

// CurrentBlockHeight returns the height of the latest block in the chain.
// The value does not necessarilly match the length of Blockchain.blocks.
func (bc *Blockchain) CurrentBlockHeight() uint32 {
	return bc.LatestBlock().Height
}

// LatestBlock returns the block at the last index of the Blockchain.blocks.
func (bc *Blockchain) LatestBlock() *block.Block {
	return bc.blocks[len(bc.blocks)-1]
}

// Difficulty returns the next mining target.
// TODO implement Difficulty
func (bc *Blockchain) Difficulty() uint32 {
	return 5
}
