package transaction

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"go_chain/common"
	"go_chain/wallet"
)

type TxType string

const (
	Coinbase TxType = "coinbase"
	Normal   TxType = "normal"
)

type Transaction struct {
	TxType
	TxHash     [32]byte
	TotalValue uint32
	Fee        uint32
	SenderAddr string
	Timestamp  int64
	Signature  *wallet.Signature
	Outs       []*Output
}

func New() *Transaction {
	return &Transaction{}
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		TxHash     string    `json:"txHash"`
		TxType     TxType    `json:"txType"`
		TotalValue uint32    `json:"totalValue"`
		Fee        uint32    `json:"fee"`
		SenderAddr string    `json:"senderAddr"`
		Timestamp  int64     `json:"timestamp"`
		Signature  string    `json:"signature"`
		Outs       []*Output `json:"outs"`
	}{
		TxHash:     fmt.Sprintf("%x", t.TxHash),
		TxType:     t.TxType,
		TotalValue: t.TotalValue,
		Fee:        t.Fee,
		SenderAddr: t.SenderAddr,
		Timestamp:  t.Timestamp,
		Signature:  t.Signature.String(),
		Outs:       t.Outs,
	})
}

func (tx *Transaction) ToBytes() []byte {
	buf := &bytes.Buffer{}
	buf.Write(tx.TxHash[:])
	buf.Write([]byte(tx.TxType))
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
	hash := sha256.New()
	tx.TxHash = sha256.Sum256(hash.Sum(tx.ToBytes()))
	return nil
}

// func (tx *Transaction) addHash(h hash.Hash, strs []string) error {
// 	for _, s := range strs {
// 		_, err := h.Write([]byte(s))
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

func (tx *Transaction) Verify() bool {
	return true
}

type Output struct {
	RecipientAddr string
	Value         uint32
}

func (o *Output) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		&struct {
			RecipientAddr string `json:"recipientAddr"`
			Value         uint32 `json:"value"`
		}{
			RecipientAddr: o.RecipientAddr,
			Value:         o.Value,
		})
}
