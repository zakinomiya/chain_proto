package transaction

import (
	"go_chain/common"
	"go_chain/wallet"
)

func NewCoinbase(w *wallet.Wallet, value uint32) *Transaction {
	tx := New()

	tx.Fee = 0
	tx.SenderAddr = ""
	tx.Timestamp = common.Timestamp()
	tx.TotalValue = value
	tx.Outs = []*Output{
		{
			Index:         0,
			RecipientAddr: w.Address(),
			Value:         value,
		},
	}
	tx.CalcHash()

	sig, _ := w.Sign(tx.TxHash[:])
	tx.Signature = sig
	return tx
}
