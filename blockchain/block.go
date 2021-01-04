package blockchain

import (
	"chain_proto/block"
	"log"
)

// Blocks returns Blockchain.blocks
func (bc *Blockchain) Blocks() []*block.Block {
	return bc.blocks
}

// ReplaceBlocks replaces the entire chain with the new one.
func (bc *Blockchain) ReplaceBlocks(blocks []*block.Block) {
	bc.blocks = blocks
}

// AddBlock adds a new block to the chain.
func (bc *Blockchain) AddBlock(block *block.Block) bool {
	if !block.Verify() {
		log.Println("info: Refused adding the block")
		return false
	}

	if err := bc.repository.Block.Insert(block); err != nil {
		log.Printf("error: failed to insert block to db. err=%v", err)
		return false
	}

	bc.blocks = append(bc.blocks, block)

	log.Printf("info: Adding new block: %x\n", block.Hash)
	log.Printf("debug: block=%+v\n", block)
	log.Printf("info: Now the length of the chain is %d:\n", bc.LatestBlock().Height)
	return true
}

func (bc *Blockchain) GetBlockByHash(hash string) (*block.Block, error) {
	b, err := bc.repository.Block.GetBlockByHash(hash)
	if err != nil {
		return nil, err
	}

	return b, nil
}
