package repository

import (
	"go_chain/account"
	"go_chain/db/models"

	"github.com/jmoiron/sqlx"
)

type AccountRepository struct {
	*database
}

func (ar *AccountRepository) GetAccount(addr string) (*account.Account, error) {
	row, err := ar.queryRow("get_account", map[string]interface{}{"addr": addr})
	if err != nil {
		return nil, err
	}

	am := &models.AccountModel{}
	if err := row.StructScan(&am); err != nil {
		return nil, err
	}

	return am.ToAccount(), nil
}

func (ar *AccountRepository) Insert(account *account.Account) error {
	return ar.insert(nil, account)
}

func (ar *AccountRepository) insert(tx *sqlx.Tx, account *account.Account) error {
	filename := "insert_account.sql"
	am := &models.AccountModel{}
	am.FromAccount(account)

	if tx != nil {
		return ar.txCommand(tx, filename, am)
	}

	return ar.command(filename, am)
}

func (ar *AccountRepository) BulkInsert(accounts []*account.Account) error {
	return ar.bulkInsert(nil, accounts)
}

func (ar *AccountRepository) bulkInsert(tx *sqlx.Tx, accounts []*account.Account) error {
	filename := "insert_account.sql"
	accountModels := []*models.AccountModel{}
	for _, account := range accounts {
		am := &models.AccountModel{}
		am.FromAccount(account)
		accountModels = append(accountModels, am)
	}

	if tx != nil {
		return ar.txCommand(tx, filename, accountModels)
	}

	return ar.command(filename, accountModels)
}
