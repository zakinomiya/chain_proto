package models

import (
	"chain_proto/account"
	"log"
)

type AccountModel struct {
	Addr    string `db:"addr"`
	Balance uint32 `db:"balance"`
}

func (am *AccountModel) FromAccount(account *account.Account) {
	log.Printf("debug: action=FromAccount. Addr=%s Balance=%d\n", account.Addr, account.Balance)
	am.Addr = account.Addr
	am.Balance = account.Balance
}

func (am *AccountModel) ToAccount() *account.Account {
	account := account.New(am.Addr)
	account.Balance = am.Balance
	return account
}
