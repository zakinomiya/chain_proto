package transaction

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"go_chain/common"

	"hash"
	"log"
)

type Transaction struct {
	TxHash     [32]byte
	TotalValue uint32
	Fee        uint32
	SenderAddr [32]byte
	Timestamp  uint64
	Outs       []*Output
}

func New() *Transaction {
	return &Transaction{}
}

func (tx *Transaction) ToBytes() []byte {
	buf := &bytes.Buffer{}
	buf.Write(tx.TxHash[:])
	buf.Write(common.IntToByteSlice(int(tx.TotalValue)))
	buf.Write(common.IntToByteSlice(int(tx.Fee)))
	buf.Write(tx.SenderAddr[:])
	buf.Write(common.IntToByteSlice(int(tx.Timestamp)))

	return buf.Bytes()
}

func (tx *Transaction) TxHashStr() string {
	return fmt.Sprintf("%x", tx.TxHash)
}

func (tx *Transaction) CalcHash() error {
	log.Println("debug: action=CalcHash")
	hash := sha256.New()

	tx.TxHash = sha256.Sum256(hash.Sum(tx.ToBytes()))
	return nil
}

func (tx *Transaction) addHash(h hash.Hash, strs []string) error {
	for _, s := range strs {
		_, err := h.Write([]byte(s))
		if err != nil {
			return err
		}
	}

	return nil
}
func (tx *Transaction) Verify() bool {
	return true
}

type Output struct {
	Index         uint32   `json:"index"`
	RecipientAddr [32]byte `json:"recipientAddr"`
	Value         uint32   `json:"value"`
	Signature     [32]byte `json:"signature"`
}

func (o *Output) Sign(privKey []byte) {
	o.Signature = [32]byte{}
}
