package transaction

import "log"

type Transaction struct {
	hash      string
	vins      []*Input
	vouts     []*Output
	timestamp uint64
	signature []byte
}

func New() *Transaction {
	return &Transaction{
		vins:  []*Input{},
		vouts: []*Output{},
	}
}

func (tx *Transaction) Sign() *Transaction {
	// TODO
	return tx
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

func (tx *Transaction) CalcHash() *Transaction {
	log.Println("action=CalcHash")
	return tx
}

type Output struct {
	index  uint32
	pubKey string
	value  uint64
}

func NewOutput(index uint32, pubKey string, value uint64) *Output {
	return &Output{index, pubKey, value}
}

type Input struct {
	index        uint32
	previousHash string
}

func NewInput(index uint32, previousHash string) *Input {
	return &Input{index, previousHash}
}
