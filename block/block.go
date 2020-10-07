package block

import (
	"crypto"
	"crypto/sha256"
	"math/rand"
	"strconv"

	"github.com/NebulousLabs/merkletree"

	"go_chain/transaction"
)

type Block struct {
	header       *BlockHeader
	hash         [32]byte
	transactions []*transaction.Transaction
	extraNonce   uint32
}

func New() *Block {
	header := NewHeader()
	return &Block{header, [32]byte{}, nil, 0}
}

func (block *Block) Header() *BlockHeader {
	return block.header
}

func (block *Block) SetExtraNonce() {
	block.extraNonce = rand.Uint32()
}

func (block *Block) Hash() [32]byte {
	return block.hash
}

func (block *Block) SetHash(hash [32]byte) {
	block.hash = hash
}

func (block *Block) SetTranscations(txs []*transaction.Transaction) {
	block.transactions = txs
}

func (block *Block) CalcMerkleTree() {
	tree := merkletree.New(crypto.SHA256.New())

	for _, tx := range block.transactions {
		h := tx.TxHash()
		tree.Push(h[:])
	}

	block.header.merkleRoot = tree.Root()
}

func (block *Block) HashBlock() [32]byte {
	sha := sha256.New()

	sha.Write([]byte(strconv.Itoa(int(block.extraNonce))))
	sha.Write(block.header.prevBlockHash)
	sha.Write(block.header.merkleRoot)
	sha.Write([]byte(strconv.Itoa(int(block.header.timestamp))))
	sha.Write([]byte(strconv.Itoa(int(block.header.bits))))
	sha.Write([]byte(strconv.Itoa(int(block.header.nonce))))

	return sha256.Sum256(sha.Sum([]byte{}))
}
