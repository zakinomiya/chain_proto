package transaction

import (
	"bytes"
	"chain_proto/common"
	"chain_proto/wallet"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/shopspring/decimal"
)

type TxType string

const (
	Coinbase TxType = "coinbase"
	Normal   TxType = "normal"
)

type Transaction struct {
	TxType
	TxHash     [32]byte
	TotalValue decimal.Decimal
	Fee        decimal.Decimal
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
		TotalValue string    `json:"totalValue"`
		Fee        string    `json:"fee"`
		SenderAddr string    `json:"senderAddr"`
		Timestamp  int64     `json:"timestamp"`
		Signature  string    `json:"signature"`
		Outs       []*Output `json:"outs"`
	}{
		TxHash:     fmt.Sprintf("%x", t.TxHash),
		TxType:     t.TxType,
		TotalValue: t.TotalValue.String(),
		Fee:        t.Fee.String(),
		SenderAddr: t.SenderAddr,
		Timestamp:  t.Timestamp,
		Signature:  t.Signature.String(),
		Outs:       t.Outs,
	})
}

func (t *Transaction) UnmarshalJSON(b []byte) error {
	tx := &struct {
		TxHash     string    `json:"txHash"`
		TxType     TxType    `json:"txType"`
		TotalValue string    `json:"totalValue"`
		Fee        string    `json:"fee"`
		SenderAddr string    `json:"senderAddr"`
		Timestamp  int64     `json:"timestamp"`
		Signature  string    `json:"signature"`
		Outs       []*Output `json:"outs"`
	}{}
	err := json.Unmarshal(b, &tx)
	if err != nil {
		return err
	}

	hash, err := hex.DecodeString(tx.TxHash)
	if err != nil {
		return err
	}

	totalValue, err := decimal.NewFromString(tx.TotalValue)
	if err != nil {
		return err
	}

	fee, err := decimal.NewFromString(tx.Fee)
	if err != nil {
		return err
	}

	t.TxHash = common.ReadByteInto32(hash)
	t.TxType = tx.TxType
	t.TotalValue = totalValue
	t.Fee = fee
	t.SenderAddr = tx.SenderAddr
	t.Timestamp = tx.Timestamp
	t.Outs = tx.Outs

	s, err := wallet.DecodeSigString(tx.Signature)
	if err != nil {
		return err
	}
	t.Signature = s
	return nil
}

func (tx *Transaction) ToBytes() []byte {
	buf := &bytes.Buffer{}
	buf.Write(tx.TxHash[:])
	buf.Write([]byte(tx.TxType))
	buf.Write([]byte(tx.TotalValue.String()))
	buf.Write([]byte(tx.Fee.String()))
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

// TODO implement
func (tx *Transaction) Verify() bool {
	return true
}

type Output struct {
	RecipientAddr string
	Value         decimal.Decimal
}

func (o *Output) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		&struct {
			RecipientAddr string `json:"recipientAddr"`
			Value         string `json:"value"`
		}{
			RecipientAddr: o.RecipientAddr,
			Value:         o.Value.String(),
		})
}
