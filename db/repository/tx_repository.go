package repository

import (
	"chain_proto/db/models"
	"chain_proto/transaction"
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
)

type TxRepository struct {
	*database
}

func (tr *TxRepository) GetByHash(hash [32]byte) (*transaction.Transaction, error) {
	row, err := tr.queryRow("get_tx_by_hash.sql", map[string]interface{}{"hash": hash[:]})
	if err != nil {
		return nil, err
	}

	txModel := &models.TxModel{}
	if err := row.StructScan(txModel); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}

		return nil, err
	}

	tx := transaction.New()
	if err := txModel.ToTx(tx); err != nil {
		return nil, err
	}

	return tx, nil
}

func (tr *TxRepository) GetByBlockHash(blockHash [32]byte) ([]*transaction.Transaction, error) {
	txModels, err := tr.getTxModelsByBlockHash(blockHash)
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

func (tr *TxRepository) getTxModelsByBlockHash(blockHash [32]byte) ([]*models.TxModel, error) {
	rows, err := tr.queryRows("get_txs_by_block_hash.sql", map[string]interface{}{"blockHash": blockHash[:]})
	if err != nil {
		return nil, err
	}

	var txModels []*models.TxModel
	for rows.Next() {
		txModel := &models.TxModel{}
		if err := rows.StructScan(txModel); err != nil {
			if err == sql.ErrNoRows {
				return nil, ErrNotFound
			}
			log.Println("error:", err)
			return nil, err
		}
		txModels = append(txModels, txModel)
	}

	return txModels, nil
}

func (tr *TxRepository) Insert(tx *transaction.Transaction) error {
	filename := "insert_tx.sql"
	txModel := &models.TxModel{}
	if err := txModel.FromTx(tx); err != nil {
		return err
	}

	return tr.command(filename, txModel)
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
