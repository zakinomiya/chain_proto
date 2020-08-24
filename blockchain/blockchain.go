package blockchain

import (
	"go_chain/block"
	"go_chain/config"
	"go_chain/repository"
	"go_chain/utils"
	"log"
	"sync"
)

type Blockchain struct {
	blocks       []*block.Block
	repositories *repository.Repositories
}

var blockchain *Blockchain
var once sync.Once

func New(conf *config.ConfigSettings) *Blockchain {
	// TODO initialise db based on config

	once.Do(func() { initializeBlockchain() })
	return blockchain
}

func initializeBlockchain() {
	blockchain = &Blockchain{
		repositories: repository.New(),
	}

	blocks, err := blockchain.repositories.BlockRepository.GetAll()

	if err != nil {
		log.Fatalln("Failed to initialise blockchain")
		return
	}

	if len(blocks) == 0 {
		log.Println("No blocks found in the db. Creating the genesis block")
		genesis := utils.NewGenesisBlock()
		blockchain.AddNewBlock(genesis)
		return
	}

	log.Println("Block record found in the db. Restoring the blockchain")
	blockchain.ReplaceBlocks(blocks)
}

func (BC *Blockchain) ServiceName() string {
	return "Blockchain"
}

func (bc *Blockchain) Start() error {
	return nil
}

func (bc *Blockchain) Stop() {
	return
}

func (bc *Blockchain) Blocks() []*block.Block {
	return bc.blocks
}

func (bc *Blockchain) ReplaceBlocks(blocks []*block.Block) {
	bc.blocks = blocks
}

func (bc *Blockchain) AddNewBlock(block *block.Block) (*Blockchain, error) {
	log.Printf("Adding new block: %#v \n", block)
	bc.blocks = append(bc.blocks, block)
	return bc, nil
}
