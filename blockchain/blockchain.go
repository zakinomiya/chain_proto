package blockchain

import (
	"go_chain/block"
	"go_chain/config"
	"log"
	"sync"
)

type Blockchain struct {
	blocks []block.Block
}

var blockchain *Blockchain
var once sync.Once

func New(conf *config.ConfigSettings) *Blockchain {
	// TODO initialise db based on config

	once.Do(func() { initializeBlockchain() })
	return blockchain
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

func initializeBlockchain() {
	blockchain = &Blockchain{[]block.Block{}}
	// if blocks == nil {
	// 	log.Println("Creating the genesis block")
	// 	genesis := utils.NewGenesisBlock()
	// 	blockchain.AddNewBlock(genesis)
	// 	return
	// }
	blockchain.ReplaceBlocks(nil)
}

func (bc *Blockchain) Blocks() []block.Block {
	return bc.blocks
}

func (bc *Blockchain) ReplaceBlocks(blocks []block.Block) {
	bc.blocks = blocks
}

func (bc *Blockchain) AddNewBlock(block *block.Block) (*Blockchain, error) {
	log.Printf("Adding new block: %#v \n", block)
	bc.blocks = append(bc.blocks, *block)
	return bc, nil
}
