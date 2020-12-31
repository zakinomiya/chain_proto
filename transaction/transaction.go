package transaction

import (
	"bytes"
	"chain_proto/common"
	"chain_proto/wallet"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
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

// TODO implements
func (t *Transaction) UnmarshalJSON(buf []byte) error {
	type te struct {
		TxHash     string    `json:"txHash"`
		TxType     TxType    `json:"txType"`
		TotalValue uint32    `json:"totalValue"`
		Fee        uint32    `json:"fee"`
		SenderAddr string    `json:"senderAddr"`
		Timestamp  int64     `json:"timestamp"`
		Signature  string    `json:"signature"`
		Outs       []*Output `json:"outs"`
	}
	teS := &te{}
	err := json.Unmarshal(buf, &teS)
	if err != nil {
		return err
	}

	bhex, err := hex.DecodeString(teS.TxHash)
	if err != nil {
		return err
	}

	t.TxHash = common.ReadByteInto32(bhex)
	t.TxType = teS.TxType
	t.TotalValue = teS.TotalValue
	t.Fee = teS.Fee
	t.SenderAddr = teS.SenderAddr
	t.Timestamp = teS.Timestamp
	t.Outs = teS.Outs

	s, err := wallet.DecodeSigString(teS.Signature)
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
