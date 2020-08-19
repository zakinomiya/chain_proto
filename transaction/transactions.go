package transaction

type Input struct {
	index uint32
	hash  string
}

type Output struct {
	index uint32
	hash  string
}

type Transaction struct {
	hash       string
	vin        []Input
	vout       []Output
	signature  []byte
	signerAddr string
}

func New() *Transaction {
	return &Transaction{}
}
