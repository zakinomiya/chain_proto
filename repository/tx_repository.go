package repository

import (
	"encoding/json"
	"go_chain/common"
	"go_chain/transaction"
	"log"
)

type txModel struct {
	txHash     []byte
	totalValue uint32
	fee        uint32
	senderAddr []byte
	timestamp  uint64
	outCount   int
	outs       []byte
}

func (r *Repository) toTx(tm *txModel, tx *transaction.Transaction) error {
	txHash, err := common.ReadByteInto32(tm.txHash)
	if err != nil {
		return err
	}

	senderAddr, err := common.ReadByteInto32(tm.senderAddr)
	if err != nil {
		return err
	}

	var outs []*transaction.Output
	err = json.Unmarshal(tm.outs, &outs)
	if err != nil || tm.outCount != len(outs) {
		return err
	}

	tx.TxHash = txHash
	tx.TotalValue = tm.totalValue
	tx.Timestamp = tm.timestamp
	tx.SenderAddr = senderAddr
	tx.Fee = tm.fee
	tx.Outs = outs

	return nil
}

func (r *Repository) fromTx(tx *transaction.Transaction, tm *txModel) error {
	outs, err := json.Marshal(tx.Outs)
	if err != nil {
		log.Println("error: failed to marshal outs to JSON format. ", err)
		return err
	}

	tm.txHash = tx.TxHash[:]
	tm.totalValue = tx.TotalValue
	tm.fee = tx.Fee
	tm.senderAddr = tx.SenderAddr[:]
	tm.timestamp = tx.Timestamp
	tm.outCount = len(tx.Outs)
	tm.outs = outs

	return nil
}

func (r *Repository) GetTxByBlockHash(blockHash [32]byte) ([]*transaction.Transaction, error) {
	txModels, err := r.getTxModelByBlockHash(blockHash)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var transactions []*transaction.Transaction
	for _, txModel := range txModels {
		var t *transaction.Transaction
		if err := r.toTx(txModel, t); err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}

	return transactions, nil
}

func (r *Repository) getTxModelByBlockHash(blockHash [32]byte) ([]*txModel, error) {
	rows, err := r.find("get_txs_by_block_hash.sql", map[string]interface{}{"blockHash": blockHash[:]})
	if err != nil {
		return nil, err
	}

	var txModels []*txModel
	for rows.Next() {
		txModel := &txModel{}
		if err := rows.StructScan(txModel); err != nil {
			return nil, err
		}
		txModels = append(txModels, txModel)
	}

	return txModels, nil
}
