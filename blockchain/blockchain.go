package blockchain

import (
	"go_chain/block"
	"go_chain/utils"
	"log"
	"sync"
)

type Blockchain struct {
	blocks []block.Block
}

var blockchain *Blockchain
var once sync.Once

func New(blocks []block.Block) *Blockchain {
	once.Do(func() { initializeBlockchain(blocks) })

	return blockchain
}

func initializeBlockchain(blocks []block.Block) {
	blockchain = &Blockchain{[]block.Block{}}
	if blocks == nil {
		log.Println("Creating the genesis block")
		genesis := utils.NewGenesisBlock()
		blockchain.AddNewBlock(genesis)
		return
	}
	blockchain.ReplaceBlocks(blocks)
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
