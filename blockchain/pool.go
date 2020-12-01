package blockchain

import "go_chain/transaction"

/// GetPooledTransactions returns transactions in the pool.
/// note this function returns the new slice of transactions.
func (bc *Blockchain) GetPooledTransactions(num int) []*transaction.Transaction {

	if len(blockchain.transactionPool) <= num {
		r := make([]*transaction.Transaction, len(bc.transactionPool))
		copy(r, bc.transactionPool)
		return r
	}

	r := make([]*transaction.Transaction, num)
	copy(r, bc.transactionPool[:num])
	return r
}

func (bc *Blockchain) deleteTxsFromPool(txs []*transaction.Transaction) {
	var newPool []*transaction.Transaction

	for _, tx := range txs {
		for _, t := range bc.transactionPool {
			if tx.TxHash == t.TxHash {
				continue
			}
			newPool = append(newPool, t)
		}
	}

	bc.transactionPool = newPool
}
