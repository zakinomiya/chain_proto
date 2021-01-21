package models

import (
	"chain_proto/account"
	"chain_proto/common"
	"chain_proto/config"
	"log"
)

const balancePrefix = "B"

type AccountModel struct {
	Addr    string `db:"addr"`
	Balance string `db:"balance"`
}

func (am *AccountModel) FromAccount(account *account.Account) {
	log.Printf("debug: action=FromAccount. Addr=%s Balance=%s\n", account.Addr, account.Balance)
	am.Addr = account.Addr
	// SQLite3 will interpret numeric strings(like "123.456") as number; so prefixing string
	am.Balance = balancePrefix + account.Balance.StringFixed(config.MaxDecimalDigit)
}

func (am *AccountModel) ToAccount() (*account.Account, error) {
	balance, err := common.ToDecimal(am.Balance, balancePrefix)
	if err != nil {
		return nil, err
	}

	account := account.New(am.Addr)
	account.Balance = balance
	return account, nil
}
