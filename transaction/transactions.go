package transaction

import (
	"bytes"
	"crypto/sha256"
	"go_chain/common"
	"hash"
	"log"
)

type Transaction struct {
	txHash        [32]byte
	totalValue    uint32
	fee           uint32
	senderAddr    [32]byte
	timestamp     uint64
	recipientAddr [32]byte
	value         uint32
	signature     [32]byte
}

func New() *Transaction {
	return &Transaction{}
}

func (tx *Transaction) ToBytes() []byte {
	buf := &bytes.Buffer{}
	buf.Write(tx.txHash[:])
	buf.Write(common.IntToByteSlice(int(tx.totalValue)))
	buf.Write(common.IntToByteSlice(int(tx.fee)))
	buf.Write(tx.senderAddr[:])
	buf.Write(common.IntToByteSlice(int(tx.timestamp)))

	return buf.Bytes()
}

func (tx *Transaction) TxHash() [32]byte {
	return tx.txHash
}

func (tx *Transaction) CalcHash() error {
	log.Println("debug: action=CalcHash")
	hash := sha256.New()

	tx.txHash = sha256.Sum256(hash.Sum(tx.ToBytes()))
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

func (tx *Transaction) Sign(privKey []byte) {
	tx.signature = [32]byte{}
}
