package transaction

import (
	"crypto/sha256"
	"hash"
	"log"
)

type Transaction struct {
	hash      [32]byte
	vins      []*Input
	vouts     []*Output
	timestamp uint64
}

func New() *Transaction {
	return &Transaction{
		vins:  []*Input{},
		vouts: []*Output{},
	}
}

func (tx *Transaction) Hash() [32]byte {
	return tx.hash
}

func (tx *Transaction) CalcHash() error {
	log.Println("action=CalcHash")
	hash := sha256.New()

	for _, vin := range tx.vins {
		if err := tx.addHash(hash, []string{string(vin.index), vin.previousHash, vin.signature}); err != nil {
			return err
		}
	}

	for _, vout := range tx.vouts {
		if err := tx.addHash(hash, []string{vout.pubKey, string(vout.value)}); err != nil {
			return err
		}
	}

	tx.hash = sha256.Sum256(hash.Sum([]byte{}))
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

func (tx *Transaction) AddInput(input *Input) *Transaction {
	log.Println("action=AddInput")
	tx.vins = append(tx.vins, input)
	return tx
}

func (tx *Transaction) AddOutput(output *Output) *Transaction {
	log.Println("action=AddOutput")
	tx.vouts = append(tx.vouts, output)
	return tx
}

type Output struct {
	pubKey string
	value  uint64
}

func NewOutput(pubKey string, value uint64) *Output {
	return &Output{pubKey, value}
}

type Input struct {
	index        uint32
	previousHash string
	signature    string
}

func NewInput(index uint32, previousHash string) *Input {
	return &Input{index, previousHash, ""}
}

func (input *Input) Sign(privKey []byte) {
	input.signature = ""
}
