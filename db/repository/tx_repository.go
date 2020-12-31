package repository

import (
	"chain_proto/db/models"
	"chain_proto/transaction"
	"log"

	"github.com/jmoiron/sqlx"
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
	rows, err := tr.queryRows("get_txs_by_block_hash.sql", map[string]interface{}{"blockHash": blockHash[:]})
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

func (tr *TxRepository) BulkInsert(txs []*transaction.Transaction) error {
	return tr.bulkInsert(nil, txs)
}

func (tr *TxRepository) bulkInsert(tx *sqlx.Tx, txs []*transaction.Transaction) error {
	filename := "insert_tx.sql"
	txModels := make([]*models.TxModel, 0, len(txs))
	for _, t := range txs {
		txModel := &models.TxModel{}
		if err := txModel.FromTx(t); err != nil {
			return err
		}
		txModels = append(txModels, txModel)
	}

	if tx != nil {
		return tr.txCommand(tx, filename, txModels)
	}

	return tr.command(filename, txModels)
}
