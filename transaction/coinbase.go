package transaction

import (
	"go_chain/wallet"
	"time"
)

func NewCoinbase(w *wallet.Wallet, value uint32) *Transaction {
	tx := New()

	tx.Fee = 0
	tx.SenderAddr = ""
	tx.Timestamp = uint64(time.Now().UnixNano())
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
