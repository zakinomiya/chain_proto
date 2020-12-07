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
	tx.TxType = "coinbase"
	sig, _ := w.Sign(append(tx.TxHash[:], w.PubKeyBytes()...))
	tx.Signature = sig
	tx.Outs = []*Output{
		{
			RecipientAddr: w.Address(),
			Value:         value,
		},
	}
	tx.CalcHash()

	return tx
}
