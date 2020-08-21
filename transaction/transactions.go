package transaction

import "log"

type Transaction struct {
	hash      string
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

func (tx *Transaction) Hash() {
	// vins[..] -> vouts[..] -> timestamp
	// input: index -> prevHash -> signature
	// output: pubKey -> value
	// []byte

	tx.hash = ""
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
