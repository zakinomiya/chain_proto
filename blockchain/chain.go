package blockchain

import (
	"fmt"
	"go_chain/block"
	"go_chain/repository"
	"go_chain/transaction"
	"log"
	"os"
	"sync"
)

type Blockchain struct {
	lock            sync.RWMutex
	chainID         uint32
	height          uint32
	blocks          []*block.Block
	transactionPool []*transaction.Transaction
	repository      *repository.Repository
}

type BlockchainInterface interface {
	CurrentBlockHeight() uint32
	Difficulty() uint32
	LatestBlock() *block.Block
	AddBlock(block *block.Block) bool
	GetPooledTransactions(num int) []*transaction.Transaction
}

var blockchain *Blockchain
var once sync.Once

func New(chainID uint32, repository *repository.Repository) *Blockchain {
	blockchain = &Blockchain{
		chainID: chainID, repository: repository, transactionPool: make([]*transaction.Transaction, 0),
	}
	return blockchain
}

func initializeBlockchain() error {
	blocks, err := blockchain.repository.GetLatestBlocks(10)

	if err != nil {
		log.Println("error: Failed to initialise blockchain")
		return err
	}

	if len(blocks) == 0 {
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
	blockchain.height = uint32(len(blocks))
	blockchain.ReplaceBlocks(blocks)
	return nil
}

func (bc *Blockchain) ServiceName() string {
	return "Blockchain"
}

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

func (bc *Blockchain) Stop() {
	return
}

func (bc *Blockchain) CurrentBlockHeight() uint32 {
	return bc.LatestBlock().Height
}

func (bc *Blockchain) LatestBlock() *block.Block {
	return bc.blocks[len(bc.blocks)-1]
}

func (bc *Blockchain) Difficulty() uint32 {
	return 5
}
