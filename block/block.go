package block

import (
	"context"
	"crypto"
	"crypto/sha256"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/NebulousLabs/merkletree"

	"go_chain/transaction"
)

type BlockHeader struct {
	prevBlockHash []byte
	merkleRoot    []byte
	timestamp     uint32
	nonce         uint32
}

func NewHeader() *BlockHeader {
	return &BlockHeader{nonce: 0}
}

type Block struct {
	mut          *sync.Mutex
	header       *BlockHeader
	hash         [32]byte
	transactions []*transaction.Transaction
	signerAddr   string
}

func New() *Block {
	header := NewHeader()
	return &Block{&sync.Mutex{}, header, [32]byte{}, nil, "Miner"}
}

func (block *Block) Hash() [32]byte {
	return block.hash
}

func (block *Block) CalcTxHash() {
	tree := merkletree.New(crypto.SHA256.New())

	for _, tx := range block.transactions {
		h := tx.Hash()
		tree.Push(h[:])
	}

	block.header.merkleRoot = tree.Root()
}

func (block *Block) HashBlock() [32]byte {
	sha := sha256.New()
	sha.Write(block.header.prevBlockHash)
	sha.Write(block.header.merkleRoot)
	sha.Write([]byte(string(block.header.timestamp)))
	sha.Write([]byte(string(block.header.nonce)))

	return sha256.Sum256(sha.Sum([]byte{}))
}

func (block *Block) AddTransaction(tx *transaction.Transaction) {
	block.transactions = append(block.transactions, tx)
}

func (block *Block) Work(difficulty int) {
	_, cancel := context.WithCancel(ctx)
	defer cancel()
	consecutiveZeros := strings.Repeat("0", difficulty)
	block.mut.Lock()
	defer block.mut.Unlock()

	for {
		blockHash := block.HashBlock()
		if strings.HasPrefix(fmt.Sprintf("%x", blockHash), consecutiveZeros) {
			log.Printf("Found valid nonce: %s", block.header.nonce)
			log.Printf("block header hash is : %x", blockHash)
			block.hash = blockHash
			ctx.Done()
		}

		block.header.nonce++
	}
}
