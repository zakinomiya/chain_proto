package transaction

import (
	"chain_proto/common"
	"chain_proto/config"
	"chain_proto/wallet"

	"github.com/shopspring/decimal"
)

func NewCoinbase(w *wallet.Wallet, value string) *Transaction {
	tx := New()

	totalValue, _ := decimal.NewFromString(value)
	tx.Fee = decimal.New(0, config.MaxDecimalDigit)
	tx.SenderAddr = ""
	tx.Timestamp = common.Timestamp()
	tx.TotalValue = totalValue
	tx.TxType = "coinbase"
	sig, _ := w.Sign(append(tx.TxHash[:], w.PubKeyBytes()...))
	tx.Signature = sig
	tx.Outs = []*Output{
		{
			RecipientAddr: w.Address(),
			Value:         totalValue,
		},
	}
	tx.CalcHash()

	return tx
}
