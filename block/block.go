package block

import (
	"go_chain/transaction"
)

type BlockHeader struct {
	prevBlockHash string
	merkleRoot    string
	timestamp     uint32
	nonce         uint32
}

func NewHeader() *BlockHeader {
	return &BlockHeader{}
}

type Block struct {
	header       *BlockHeader
	hash         string
	transactions []*transaction.Transaction
	signerAddr   string
}

func New() *Block {
	header := NewHeader()
	return &Block{header, "some hash", nil, "Miner"}
}

func (block *Block) Hash() string {
	return block.hash
}

func (block *Block) SetHash(hash string) *Block {
	block.hash = hash
	return block
}

func (block *Block) HashBlock() {

}
