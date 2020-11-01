package blockchain

import (
	"fmt"
	"go_chain/block"
	"go_chain/transaction"
	"log"
	"strings"
)

func (bc *Blockchain) Blocks() []*block.Block {
	return bc.blocks
}

func (bc *Blockchain) ReplaceBlocks(blocks []*block.Block) {
	bc.blocks = blocks
}

func (bc *Blockchain) AddBlock(block *block.Block) bool {
	// if !bc.verifyBlock(block) {
	// 	log.Println("info: Refused adding the block")
	// 	return false
	// }
	// log.Printf("debug: Adding new block: %#v \n", block)
	// log.Printf("debug: Now the length of the chain is %d:\n", len(bc.blocks))
	// bc.blocks = append(bc.blocks, blocke
	bc.repository.Insert(block)
	return true
}

/// Check if reveived blocks are valid blocks
/// Criteria:
///  - Block hash is valid
///     - hash = the calculated result from properties in the block
///     - hash clears its difficulty
func (bc *Blockchain) verifyBlock(block *block.Block) bool {
	hash := block.Hash

	if hash != block.HashBlock() {
		log.Printf("info: block hash is invalid. presented: %x, calculated: %x\n", hash, block.HashBlock())
		return false
	}

	zeros := strings.Repeat("0", int(block.Bits))

	if !strings.HasPrefix(fmt.Sprintf("%x", hash), zeros) {
		log.Printf("info: block hash does not meet the difficulty. difficulty %d, presented hash %x\n ", block.Bits, hash)
		return false
	}

	return true
}

func (bc *Blockchain) GenerateBlock(txs []*transaction.Transaction) *block.Block {
	block := block.New()
	block.SetExtranNonce()
	block.CalcMerkleTree()
	return block
}
