package outer

import (
	"go_chain/block"
	"go_chain/transaction"
)

type Outer struct{}

func New() *Outer {
	return &Outer{}
}

func (o *Outer) SendBlock(block *block.Block) (string, error) {
	return "", nil
}

func (o *Outer) SendTransaction(tx *transaction.Transaction) (string, error) {
	return "", nil
}

func (o *Outer) SyncBlock() error {
	return nil
}
