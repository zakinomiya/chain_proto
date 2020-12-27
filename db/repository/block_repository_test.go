package repository

import (
	"go_chain/block"
	"go_chain/db/models"
	"testing"
)

func TestProcessTxs(t *testing.T) {
	br := &BlockRepository{}
	// txNum := 10

	// w, _ := wallet.New()
	// txs := []*transaction.Transaction{}
	// txs = append(txs, transaction.NewCoinbase(w, 10))

	// for i := 0; i < txNum-1; i++ {
	// 	tx := transaction.New()
	// 	tx.SenderAddr = fmt.Sprintf("sender-%d", i)
	// 	out := &transaction.Output{}
	// 	out.RecipientAddr = fmt.Sprintf("recipient-%d", i)
	// 	out.Value = 10
	// 	tx.Outs = []*transaction.Output{out}
	// 	txs = append(txs, tx)
	// }

	gen, _ := block.NewGenesisBlock()
	t.Log(gen.Transactions)
	accounts, err := br.processTxs(gen.Transactions)
	if err != nil {
		t.Error(err)
	}

	t.Log(accounts)

	am := &models.AccountModel{}
	am.FromAccount(accounts[0])
	t.Log(am)
}
