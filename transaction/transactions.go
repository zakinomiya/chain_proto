package transaction

import (
	"crypto/sha256"
	"hash"
	"log"
)

type Transaction struct {
	txHash     [32]byte
	totalValue uint32
	fee        uint32
	senderAddr [32]byte
	outs       []*Output
	timestamp  uint64
}

func New() *Transaction {
	return &Transaction{}
}

func (tx *Transaction) TxHash() [32]byte {
	return tx.txHash
}

func (tx *Transaction) CalcHash() error {
	log.Println("debug: action=CalcHash")
	hash := sha256.New()

	for _, o := range tx.outs {
		if err := tx.addHash(hash, []string{string(o.value)}); err != nil {
			return err
		}
	}

	tx.txHash = sha256.Sum256(hash.Sum([]byte{}))
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

func (tx *Transaction) AddOutput(output *Output) *Transaction {
	log.Println("debug: action=AddOutput")
	tx.outs = append(tx.outs, output)
	return tx
}

func (tx *Transaction) Verify() bool {
	return true
}

type Output struct {
	recipientAddr [32]byte
	value         uint32
	signature     [32]byte
}

func NewOutput() *Output {
	return &Output{}
}

func (o *Output) Sign(privKey []byte) {
	o.signature = [32]byte{}
}
