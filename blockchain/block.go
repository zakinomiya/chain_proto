package blockchain

import (
	"go_chain/block"
	"log"
)

func (bc *Blockchain) Blocks() []*block.Block {
	return bc.blocks
}

func (bc *Blockchain) ReplaceBlocks(blocks []*block.Block) {
	bc.blocks = blocks
}

func (bc *Blockchain) AddBlock(block *block.Block) bool {
	if !block.Verify() {
		log.Println("info: Refused adding the block")
		return false
	}
	log.Printf("info: Adding new block: %x\n", block.Hash)
	log.Printf("debug: block=%+v\n", block)
	log.Printf("info: Now the length of the chain is %d:\n", len(bc.blocks))
	bc.blocks = append(bc.blocks, block)
	bc.repository.Insert(block)
	return true
}
