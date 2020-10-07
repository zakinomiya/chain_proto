package block

import "time"

type BlockHeader struct {
	prevBlockHash []byte
	merkleRoot    []byte
	timestamp     uint32
	bits          uint32
	nonce         uint32
}

func NewHeader() *BlockHeader {
	return &BlockHeader{nonce: 0, timestamp: uint32(time.Now().UTC().Unix())}
}

func (block *Block) Bits() uint32 {
	return block.header.bits
}

func (block *Block) Timestamp() uint32 {
	return block.header.timestamp
}

func (block *Block) MerkleRoot() []byte {
	return block.header.merkleRoot
}

func (block *Block) IncrementNonce() {
	block.header.nonce += 1
}

func (block *Block) Nonce() uint32 {
	return block.header.nonce
}
