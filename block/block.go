package block

import (
	"bytes"
	"crypto"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"go_chain/common"
	"go_chain/transaction"
	"log"
	"math/rand"
	"strconv"
	"strings"

	"github.com/NebulousLabs/merkletree"
)

type BlockHeader struct {
	PrevBlockHash [32]byte `json:"prevBlockHash"`
	MerkleRoot    []byte   `json:"merkleRoot"`
	Timestamp     int64    `json:"timestamp"`
	Bits          uint32   `json:"bits"`
	Nonce         uint32   `json:"nonce"`
}

func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Height        uint32                     `json:"height"`
		Hash          string                     `json:"hash"`
		Transactions  []*transaction.Transaction `json:"transactions"`
		ExtraNonce    uint32                     `json:"extraNonce"`
		PrevBlockHash string                     `json:"prevBlockHash"`
		MerkleRoot    string                     `json:"merkleRoot"`
		Timestamp     int64                      `json:"timestamp"`
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
	b := &Block{
		Height:       height,
		Transactions: txs,
		BlockHeader: &BlockHeader{
			Bits:          bits,
			PrevBlockHash: prevBlockHash,
			Timestamp:     common.Timestamp(),
			Nonce:         0,
		},
	}
	b.SetExtranNonce()
	b.SetMerkleRoot()

	return b
}

func (b *Block) SetExtranNonce() {
	b.ExtraNonce = rand.Uint32()
}

func (b *Block) SetMerkleRoot() {
	b.MerkleRoot = b.calcMerkleRoot()
}

func (b *Block) calcMerkleRoot() []byte {
	tree := merkletree.New(crypto.SHA256.New())

	for _, tx := range b.Transactions {
		h := tx.TxHash
		tree.Push(h[:])
	}

	return tree.Root()
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

/// Verify checks if received blocks are valid blocks
/// Criteria:
///  - Block hash is valid
///     - hash = the calculated result from properties in the block
///     - hash clears its difficulty
func (b *Block) Verify() bool {
	coinbaseCount := 0
	for _, tx := range b.Transactions {
		if tx.TxType == "coinbase" {
			if coinbaseCount == 1 {
				log.Println("error: this block contains more than 2 coinbase transactions.")
				return false
			}
			coinbaseCount++
		}

		if !tx.Verify() {
			log.Println("error: block contains a invalid transaction.")
			return false
		}
	}

	if bytes.Compare(b.MerkleRoot, b.calcMerkleRoot()) != 0 {
		log.Println("info: failed to verify the merkle root.")
		return false
	}

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
