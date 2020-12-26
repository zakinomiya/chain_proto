package models

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
	SenderAddr string
	Timestamp  int64
	OutCount   int
	Outs       []byte
}

func (tm *TxModel) ToTx(tx *transaction.Transaction) error {

	var outs []*transaction.Output
	err := json.Unmarshal(tm.Outs, &outs)
	if err != nil || tm.OutCount != len(outs) {
		return err
	}

	tx.TxHash = common.ReadByteInto32(tm.TxHash)
	tx.TotalValue = tm.TotalValue
	tx.Timestamp = tm.Timestamp
	tx.SenderAddr = tm.SenderAddr
	tx.Fee = tm.Fee
	tx.Outs = outs

	return nil
}

func (tm *TxModel) FromTx(tx *transaction.Transaction) error {
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
