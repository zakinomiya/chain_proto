package transaction

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"go_chain/common"
	"go_chain/wallet"

	"hash"
	"log"
)

type Transaction struct {
	TxHash     [32]byte
	TotalValue uint32
	Fee        uint32
	SenderAddr string
	Timestamp  uint64
	Outs       []*Output
	Signature  *wallet.Signature
}

func New() *Transaction {
	return &Transaction{}
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		TxHash     string    `json:"txHash"`
		TotalValue uint32    `json:"totalValue"`
		Fee        uint32    `json:"fee"`
		SenderAddr string    `json:"senderAddr"`
		Timestamp  uint64    `json:"timestamp"`
		Outs       []*Output `json:"outs"`
		Signature  string    `json:"signature"`
	}{
		TxHash:     fmt.Sprintf("%x", t.TxHash),
		TotalValue: t.TotalValue,
		Fee:        t.Fee,
		SenderAddr: t.SenderAddr,
		Timestamp:  t.Timestamp,
		Outs:       t.Outs,
		Signature:  t.Signature.String(),
	})
}

func (tx *Transaction) ToBytes() []byte {
	buf := &bytes.Buffer{}
	buf.Write(tx.TxHash[:])
	buf.Write(common.IntToByteSlice(int(tx.TotalValue)))
	buf.Write(common.IntToByteSlice(int(tx.Fee)))
	buf.Write([]byte(tx.SenderAddr))
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
	Index         uint32
	RecipientAddr string
	Value         uint32
}

func (o *Output) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		&struct {
			Index         uint32 `json:"index"`
			RecipientAddr string `json:"recipientAddr"`
			Value         uint32 `json:"value"`
		}{
			Index:         o.Index,
			RecipientAddr: o.RecipientAddr,
			Value:         o.Value,
		})
}
