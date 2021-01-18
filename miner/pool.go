package miner

import "chain_proto/transaction"

/// GetPooledTransactions returns transactions in the pool.
/// note this function returns the new slice of transactions.
func (m *Miner) GetPooledTransactions(num int) []*transaction.Transaction {
	if len(m.transactionPool) <= num {
		r := make([]*transaction.Transaction, len(m.transactionPool))
		copy(r, m.transactionPool)
		return r
	}

	r := make([]*transaction.Transaction, num)
	copy(r, m.transactionPool[:num])
	return r
}

func (m *Miner) AddTransaction(tx *transaction.Transaction) {
	m.transactionPool = append(m.transactionPool, tx)
}

func (m *Miner) deleteTxsFromPool(txs []*transaction.Transaction) {
	var newPool []*transaction.Transaction

	for _, tx := range txs {
		for _, t := range m.transactionPool {
			if tx.TxHash == t.TxHash {
				continue
			}
			newPool = append(newPool, t)
		}
	}

	m.transactionPool = newPool
}
