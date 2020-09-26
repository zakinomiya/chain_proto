package block

import (
	"crypto"
	"crypto/sha256"
	"math/rand"
	"strconv"
	"sync"
	"time"

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
	return &BlockHeader{nonce: 0, timestamp: uint32(time.Now().UTC().Unix())}
}

type Block struct {
	mut          *sync.Mutex
	header       *BlockHeader
	hash         [32]byte
	transactions []*transaction.Transaction
	signerAddr   string
	extraNonce   uint32
}

func New() *Block {
	header := NewHeader()
	return &Block{&sync.Mutex{}, header, [32]byte{}, nil, "Miner", 0}
}

func (block *Block) Header() *BlockHeader {
	return block.header
}

func (block *Block) IncrementNonce() {
	block.header.nonce += 1
}

func (block *Block) Nonce() uint32 {
	return block.header.nonce
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

	sha.Write([]byte(strconv.Itoa(int(block.extraNonce))))
	sha.Write(block.header.prevBlockHash)
	sha.Write(block.header.merkleRoot)
	sha.Write([]byte(strconv.Itoa(int(block.header.timestamp))))
	sha.Write([]byte(strconv.Itoa(int(block.header.nonce))))

	return sha256.Sum256(sha.Sum([]byte{}))
}
