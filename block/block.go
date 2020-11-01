package block

import (
	"crypto"
	"crypto/sha256"
	"go_chain/transaction"
	"math/rand"
	"strconv"
	"time"

	"github.com/NebulousLabs/merkletree"
)

type BlockHeader struct {
	PrevBlockHash [32]byte `db:"prevBlockHash"`
	MerkleRoot    []byte   `db:"merkleRoot"`
	Timestamp     uint32   `db:"timestamp"`
	Bits          uint32   `db:"bits"`
	Nonce         uint32   `db:"nonce"`
}

func NewHeader() *BlockHeader {
	return &BlockHeader{Nonce: 0, Timestamp: uint32(time.Now().UTC().Unix())}
}

func (block *Block) IncrementNonce() {
	block.Nonce += 1
}

type Block struct {
	Height       uint32                     `db:"height"`
	Hash         [32]byte                   `db:"hash"`
	Transactions []*transaction.Transaction `db:"transactions"`
	ExtraNonce   uint32                     `db:"extraNonce"`
	*BlockHeader
}

func New() *Block {
	header := NewHeader()
	return &Block{BlockHeader: header}
}

func (block *Block) SetExtranNonce() {
	block.ExtraNonce = rand.Uint32()
}

func (block *Block) CalcMerkleTree() {
	tree := merkletree.New(crypto.SHA256.New())

	// for _, tx := range block.Transactions {
	// 	h := tx.TxHash()
	// 	tree.Push(h[:])
	// }

	block.BlockHeader.MerkleRoot = tree.Root()
}

func (block *Block) HashBlock() [32]byte {
	sha := sha256.New()

	sha.Write([]byte(strconv.Itoa(int(block.ExtraNonce))))
	sha.Write(block.PrevBlockHash[:])
	sha.Write([]byte(block.BlockHeader.MerkleRoot))
	sha.Write([]byte(strconv.Itoa(int(block.BlockHeader.Timestamp))))
	sha.Write([]byte(strconv.Itoa(int(block.BlockHeader.Bits))))
	sha.Write([]byte(strconv.Itoa(int(block.BlockHeader.Nonce))))

	return sha256.Sum256(sha.Sum([]byte{}))
}
