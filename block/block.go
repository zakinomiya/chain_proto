package block

import (
	"crypto"
	"crypto/sha256"
	"encoding/json"
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
	PrevBlockHash [32]byte `json:"prevBlockHash"`
	MerkleRoot    []byte   `json:"merkleRoot"`
	Timestamp     uint32   `json:"timestamp"`
	Bits          uint32   `json:"bits"`
	Nonce         uint32   `json:"nonce"`
}

func NewHeader() *BlockHeader {
	return &BlockHeader{Nonce: 0, Timestamp: uint32(time.Now().UTC().Unix())}
}

func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Height        uint32                     `json:"height"`
		Hash          string                     `json:"hash"`
		Transactions  []*transaction.Transaction `json:"transactions"`
		ExtraNonce    uint32                     `json:"extraNonce"`
		PrevBlockHash string                     `json:"prevBlockHash"`
		MerkleRoot    string                     `json:"merkleRoot"`
		Timestamp     uint32                     `json:"timestamp"`
		Bits          uint32                     `json:"bits"`
		Nonce         uint32                     `json:"nonce"`
	}{
		Height:        b.Height,
		Hash:          fmt.Sprintf("%x", b.Hash),
		Transactions:  b.Transactions,
		ExtraNonce:    b.ExtraNonce,
		PrevBlockHash: fmt.Sprintf("%x", b.PrevBlockHash),
		MerkleRoot:    fmt.Sprintf("%x", b.MerkleRoot),
		Timestamp:     b.Timestamp,
		Bits:          b.Bits,
		Nonce:         b.Nonce,
	})
}

func (b *Block) IncrementNonce() {
	b.Nonce += 1
}

type Block struct {
	Height       uint32                     `json:"height"`
	Hash         [32]byte                   `json:"hash"`
	Transactions []*transaction.Transaction `json:"transactions"`
	ExtraNonce   uint32                     `json:"extraNonce"`
	*BlockHeader
}

func New(height uint32, bits uint32, prevBlockHash [32]byte, txs []*transaction.Transaction) *Block {
	b := &Block{}
	b.Transactions = txs
	b.Bits = bits
	b.PrevBlockHash = prevBlockHash
	b.Timestamp = uint32(time.Now().Unix())
	b.Nonce = 0
	b.SetExtranNonce()
	b.CalcMerkleTree()

	return b
}

func (b *Block) SetExtranNonce() {
	b.ExtraNonce = rand.Uint32()
}

func (b *Block) CalcMerkleTree() {
	tree := merkletree.New(crypto.SHA256.New())

	for _, tx := range b.Transactions {
		h := tx.TxHash
		tree.Push(h[:])
	}

	b.BlockHeader.MerkleRoot = tree.Root()
}

func (b *Block) HashBlock() [32]byte {
	sha := sha256.New()

	sha.Write([]byte(strconv.Itoa(int(b.ExtraNonce))))
	sha.Write(b.PrevBlockHash[:])
	sha.Write([]byte(b.BlockHeader.MerkleRoot))
	sha.Write([]byte(strconv.Itoa(int(b.BlockHeader.Timestamp))))
	sha.Write([]byte(strconv.Itoa(int(b.BlockHeader.Bits))))
	sha.Write([]byte(strconv.Itoa(int(b.BlockHeader.Nonce))))

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
