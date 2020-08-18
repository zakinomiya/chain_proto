package blockchain

import (
	"go_chain/block"
	"go_chain/utils"
	"log"
)

type Blockchain struct {
	blocks []block.Block
}

func New(blocks []block.Block) *Blockchain {
	if blocks == nil {
		log.Println("Creating the genesis block")
		genesis := utils.NewGenesisBlock()
		return &Blockchain{[]block.Block{*genesis}}
	}

	return &Blockchain{blocks}
}

func (bc *Blockchain) Blocks() []block.Block {
	return bc.blocks
}

func (bc *Blockchain) AddNewBlock(block *block.Block) (*Blockchain, error) {
	log.Printf("Adding new block: %#v \n", block)
	bc.blocks = append(bc.blocks, *block)
	return bc, nil
}
