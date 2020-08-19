package block

import (
	"go_chain/transaction"
	"time"
)

type Block struct {
	timestamp    int64
	hash         string
	amount       int
	transactions []transaction.Transaction
	signerAddr   string
}

func New() *Block {
	return &Block{time.Now().Unix(), "some hash", 100, []transaction.Transaction{}, "Miner"}
}

func (block *Block) Hash() string {
	return block.hash
}

func (block *Block) Amount() int {
	return block.amount
}

func (block *Block) SetHash(hash string) *Block {
	block.hash = hash
	return block
}

func (block *Block) SetAmount(amount int) *Block {
	block.amount = amount
	return block
}
