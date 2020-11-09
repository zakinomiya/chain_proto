package block

import (
	"crypto"
	"crypto/sha256"
	"fmt"
	"go_chain/transaction"
	"log"
	"math/rand"
	"strconv"
	"strings"
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

func New(height uint32, bits uint32, prevBlockHash [32]byte, txs []*transaction.Transaction) *Block {
	block := &Block{}
	block.Transactions = txs
	block.Bits = bits
	block.PrevBlockHash = prevBlockHash
	block.Timestamp = uint32(time.Now().Unix())
	block.Nonce = 0
	block.SetExtranNonce()
	block.CalcMerkleTree()

	return block
}

func (block *Block) SetExtranNonce() {
	block.ExtraNonce = rand.Uint32()
}

func (block *Block) CalcMerkleTree() {
	tree := merkletree.New(crypto.SHA256.New())

	for _, tx := range block.Transactions {
		h := tx.TxHash
		tree.Push(h[:])
	}

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

func (b *Block) TxCount() int {
	return len(b.Transactions)
}

/// Check if reveived blocks are valid blocks
/// Criteria:
///  - Block hash is valid
///     - hash = the calculated result from properties in the block
///     - hash clears its difficulty
func (b *Block) Verify() bool {

	if b.Hash != b.HashBlock() {
		log.Printf("info: block hash is invalid. presented: %x, calculated: %x\n", b.Hash, b.HashBlock())
		return false
	}

	zeros := strings.Repeat("0", int(b.Bits))

	if !strings.HasPrefix(fmt.Sprintf("%x", b.Hash), zeros) {
		log.Printf("info: block hash does not meet the difficulty. difficulty %d, presented hash %x\n ", b.Bits, b.Hash)
		return false
	}

	return true
}
