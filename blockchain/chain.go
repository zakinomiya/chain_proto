package blockchain

import (
	"fmt"
	"go_chain/block"
	"go_chain/repository"
	"go_chain/transaction"
	"go_chain/utils"
	"log"
	"os"
	"sync"
)

type BlockchainInterface interface {
	Height() uint32
	GenerateBlock(txs []*transaction.Transaction) *block.Block
	AddBlock(block *block.Block) bool
}

type Blockchain struct {
	lock       sync.RWMutex
	chainID    uint32
	height     uint32
	blocks     []*block.Block
	repository *repository.Repository
}

var blockchain *Blockchain
var once sync.Once

func New(chainID uint32, repository *repository.Repository) *Blockchain {
	blockchain = &Blockchain{
		chainID: chainID, repository: repository,
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
		genesis := utils.NewGenesisBlock()
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

func (bc *Blockchain) Height() uint32 {
	return bc.height
}
