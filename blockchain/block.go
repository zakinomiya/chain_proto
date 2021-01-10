package blockchain

import (
	"chain_proto/block"
	"chain_proto/db/repository"
	"errors"
	"fmt"
	"log"
)

// Blocks returns Blockchain.blocks
func (bc *Blockchain) Blocks() []*block.Block {
	return bc.blocks
}

func (bc *Blockchain) GetBlocks(offset uint32, limit uint32) ([]*block.Block, error) {
	blks, err := bc.repository.Block.GetByRange(offset, limit)
	if err != nil {
		if err == repository.ErrNotFound {
			return make([]*block.Block, 0), err
		}
		return nil, err
	}

	return blks, nil
}

// AddBlock adds a new block to the chain.
func (bc *Blockchain) AddBlock(block *block.Block) error {
	if !block.Verify() {
		return errors.New("error: failed in block verification")
	}

	if err := bc.repository.Block.Insert(block); err != nil {
		return fmt.Errorf("error: failed to insert block to db. err=%v", err)
	}

	bc.blocks = append(bc.blocks, block)

	log.Printf("info: Adding new block: %x\n", block.Hash)
	log.Printf("debug: block=%+v\n", block)
	log.Printf("info: Now the length of the chain is %d:\n", bc.LatestBlock().Height)
	return nil
}

// GetBlockByHash returns a block with the given block hash
func (bc *Blockchain) GetBlockByHash(hash [32]byte) (*block.Block, error) {
	return bc.repository.Block.GetByHash(hash)
}

// GetBlockByHeight returns a block with the given block hash
func (bc *Blockchain) GetBlockByHeight(height uint32) (*block.Block, error) {
	return bc.repository.Block.GetByHeight(height)
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

// ReplaceBlocks replaces the entire chain with the new one.
func (bc *Blockchain) replaceBlocks(blocks []*block.Block) {
	bc.blocks = blocks
}
