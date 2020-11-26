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
	sig, _ := w.Sign(append(tx.TxHash[:], w.PubKeyBytes()...))
	tx.Outs = []*Output{
		{
			Index:     0,
			Signature: sig,
			Value:     value,
		},
	}
	tx.CalcHash()

	return tx
}
