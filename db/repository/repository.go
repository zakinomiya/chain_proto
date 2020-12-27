package repository

import (
	"go_chain/account"
	"go_chain/block"
	"go_chain/transaction"
	"io/ioutil"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Repository struct {
	Block interface {
		GetBlockByHash(hash string) (*block.Block, error)
		GetBlocksByRange(start uint32, end uint32) ([]*block.Block, error)
		GetLatestBlock() (*block.Block, error)
		Insert(b *block.Block) error
	}
	Tx interface {
		GetTxsByBlockHash(blockHash [32]byte) ([]*transaction.Transaction, error)
		BulkInsert(txs []*transaction.Transaction) error
	}
	Account interface {
		BulkInsert(accounts []*account.Account) error
		Insert(account *account.Account) error
		GetAccount(addr string) (*account.Account, error)
	}
}

func New(dbPath string, dbDriver string) (*Repository, error) {
	log.Printf("debug: DBInfo driver=%s, dbpath=%s\n", dbDriver, dbPath)
	db, err := sqlx.Open(dbDriver, dbPath)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	database := &database{
		db:       db,
		dbPath:   dbPath,
		dbDriver: dbDriver,
	}

	sql, err := ioutil.ReadFile("./db/sql/initialize.sql")
	if err != nil {
		return nil, err
	}

	_, err = database.db.Exec(string(sql))
	if err != nil {
		return nil, err
	}

	ar := &AccountRepository{database: database}
	tr := &TxRepository{database: database}

	return &Repository{
		Block: &BlockRepository{
			database:    database,
			account:     ar,
			transcation: tr,
		},
		Tx:      tr,
		Account: ar,
	}, nil
}
