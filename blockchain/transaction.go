package blockchain

import (
	"chain_proto/db/repository"
	"chain_proto/transaction"
	"errors"
)

// GetTxsByBlockHash returns a set of transactions in a block.
// When no block
func (bc *Blockchain) GetTxsByBlockHash(blockHash [32]byte) ([]*transaction.Transaction, error) {
	txs, err := bc.r.Tx.GetByBlockHash(blockHash)
	if err != nil {
		if err == repository.ErrNotFound {
			return make([]*transaction.Transaction, 0), nil
		}
		return nil, err
	}

	return txs, nil
}

// GetTransactionByHash returns a transaction with the given transaction hash.
func (bc *Blockchain) GetTransactionByHash(hash [32]byte) (*transaction.Transaction, error) {
	tx, err := bc.r.Tx.GetByHash(hash)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// AddTransaction will pass a transaction to its registered miner.
// Calling this function does not garantee the transaction will be successfully stored in block.
func (bc *Blockchain) AddTransaction(tx *transaction.Transaction) error {
	if !tx.Verify() {
		return errors.New("invalid transaction")
	}

	bc.m.AddTransaction(tx)
	return nil
}
