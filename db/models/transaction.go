package models

import (
	"chain_proto/common"
	"chain_proto/transaction"
	"chain_proto/wallet"
	"encoding/json"
	"log"
)

const (
	txPrefix = "T"
)

type TxModel struct {
	TxHash []byte `db:"txHash"`
	TxType string `db:"txType"`
	//	BlockHash  string `db:"lockHash"`
	TotalValue string `db:"totalValue"`
	Fee        string `db:"fee"`
	SenderAddr string `db:"senderAddr"`
	Timestamp  int64  `db:"timestamp"`
	Signature  string `db:"signature"`
	OutCount   int    `db:"outCount"`
	Outs       []byte `db:"outs"`
}

func (tm *TxModel) ToTx(tx *transaction.Transaction) error {
	var outs []*transaction.Output
	err := json.Unmarshal(tm.Outs, &outs)
	if err != nil || tm.OutCount != len(outs) {
		return err
	}

	totalValue, err := common.ToDecimal(tm.TotalValue, txPrefix)
	if err != nil {
		return err
	}

	fee, err := common.ToDecimal(tm.Fee, txPrefix)
	if err != nil {
		return err
	}

	sig, err := wallet.DecodeSigString(tm.Signature)
	if err != nil {
		return err
	}

	tx.TxHash = common.ReadByteInto32(tm.TxHash)
	tx.TxType = transaction.TxType(tm.TxType)
	tx.TotalValue = totalValue
	tx.Timestamp = tm.Timestamp
	tx.SenderAddr = tm.SenderAddr
	tx.Signature = sig
	tx.Fee = fee
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
	tm.TxType = string(tx.TxType)
	tm.TotalValue = txPrefix + tx.TotalValue.String()
	tm.Fee = txPrefix + tx.Fee.String()
	tm.SenderAddr = tx.SenderAddr[:]
	tm.Timestamp = tx.Timestamp
	tm.Signature = tx.Signature.String()
	tm.OutCount = len(tx.Outs)
	tm.Outs = outs

	return nil
}
