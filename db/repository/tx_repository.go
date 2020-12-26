package repository

import (
	"go_chain/db/models"
	"go_chain/transaction"
	"log"
)

type TxRepository struct {
	*database
}

func (tr *TxRepository) GetTxsByBlockHash(blockHash [32]byte) ([]*transaction.Transaction, error) {
	txModels, err := tr.getTxModelByBlockHash(blockHash)
	if err != nil {
		log.Println("error:", err)
		return nil, err
	}

	var transactions []*transaction.Transaction
	for _, txModel := range txModels {
		var t *transaction.Transaction
		if err := txModel.ToTx(t); err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}

	return transactions, nil
}

func (tr *TxRepository) getTxModelByBlockHash(blockHash [32]byte) ([]*models.TxModel, error) {
	rows, err := tr.query("get_txs_by_block_hash.sql", map[string]interface{}{"blockHash": blockHash[:]})
	if err != nil {
		return nil, err
	}

	var txModels []*models.TxModel
	for rows.Next() {
		txModel := &models.TxModel{}
		if err := rows.StructScan(txModel); err != nil {
			log.Println("error:", err)
			return nil, err
		}
		txModels = append(txModels, txModel)
	}

	return txModels, nil
}
