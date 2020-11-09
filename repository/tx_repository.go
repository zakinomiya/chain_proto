package repository

import (
	"encoding/json"
	"go_chain/common"
	"go_chain/transaction"
	"log"
)

type TxModel struct {
	TxHash     []byte
	TotalValue uint32
	Fee        uint32
	SenderAddr []byte
	Timestamp  uint64
	OutCount   int
	Outs       []byte
}

func (r *Repository) toTx(tm *TxModel, tx *transaction.Transaction) error {
	txHash, err := common.ReadByteInto32(tm.TxHash)
	if err != nil {
		return err
	}

	senderAddr, err := common.ReadByteInto32(tm.SenderAddr)
	if err != nil {
		return err
	}

	var outs []*transaction.Output
	err = json.Unmarshal(tm.Outs, &outs)
	if err != nil || tm.OutCount != len(outs) {
		return err
	}

	tx.TxHash = txHash
	tx.TotalValue = tm.TotalValue
	tx.Timestamp = tm.Timestamp
	tx.SenderAddr = senderAddr
	tx.Fee = tm.Fee
	tx.Outs = outs

	return nil
}

func (r *Repository) fromTx(tx *transaction.Transaction, tm *TxModel) error {
	outs, err := json.Marshal(tx.Outs)
	if err != nil {
		log.Println("error: failed to marshal outs to JSON format. ", err)
		return err
	}

	tm.TxHash = tx.TxHash[:]
	tm.TotalValue = tx.TotalValue
	tm.Fee = tx.Fee
	tm.SenderAddr = tx.SenderAddr[:]
	tm.Timestamp = tx.Timestamp
	tm.OutCount = len(tx.Outs)
	tm.Outs = outs

	return nil
}

func (r *Repository) GetTxByBlockHash(blockHash [32]byte) ([]*transaction.Transaction, error) {
	TxModels, err := r.getTxModelByBlockHash(blockHash)
	if err != nil {
		log.Println("error:", err)
		return nil, err
	}

	var transactions []*transaction.Transaction
	for _, TxModel := range TxModels {
		var t *transaction.Transaction
		if err := r.toTx(TxModel, t); err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}

	return transactions, nil
}

func (r *Repository) getTxModelByBlockHash(blockHash [32]byte) ([]*TxModel, error) {
	rows, err := r.find("get_txs_by_block_hash.sql", map[string]interface{}{"blockHash": blockHash[:]})
	if err != nil {
		return nil, err
	}

	var txModels []*TxModel
	for rows.Next() {
		txModel := &TxModel{}
		if err := rows.StructScan(txModel); err != nil {
			log.Println("error:", err)
			return nil, err
		}
		txModels = append(txModels, txModel)
	}

	return txModels, nil
}
