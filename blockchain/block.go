package blockchain

import (
	"chain_proto/block"
	"chain_proto/db/repository"
	"context"
	"errors"
	"fmt"
	"log"
)

// Blocks returns Blockchain.blocks
func (bc *Blockchain) Blocks() []*block.Block {
	return bc.blocks
}

// GetBlocks returns blocks at block height within the given offset and offset+limit.
func (bc *Blockchain) GetBlocks(offset uint32, limit uint32) ([]*block.Block, error) {
	blks, err := bc.r.Block.GetByRange(offset, limit)
	if err != nil {
		if err == repository.ErrNotFound {
			return make([]*block.Block, 0), err
		}
		return nil, err
	}

	return blks, nil
}

// AddBlock adds a new block to the chain.
func (bc *Blockchain) AddBlock(b *block.Block) error {
	if !b.Verify() {
		return errors.New("error: failed in block verification")
	}

	if err := bc.r.Block.Insert(b); err != nil {
		return fmt.Errorf("error: failed to insert block to db. err=%v", err)
	}

	bc.c.PropagateBlock(context.Background(), b)
	bc.blocks = append(bc.blocks, b)

	log.Printf("info: Adding new block: %x\n", b.Hash)
	log.Printf("debug: block=%+v\n", b)
	log.Printf("info: Now the length of the chain is %d:\n", bc.LatestBlock().Height)
	return nil
}

// GetBlockByHash returns a block with the given block hash
func (bc *Blockchain) GetBlockByHash(hash [32]byte) (*block.Block, error) {
	return bc.r.Block.GetByHash(hash)
}

// GetBlockByHeight returns a block with the given block hash
func (bc *Blockchain) GetBlockByHeight(height uint32) (*block.Block, error) {
	return bc.r.Block.GetByHeight(height)
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
