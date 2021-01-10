package repository

import (
	"chain_proto/account"
	"chain_proto/block"
	"chain_proto/config"
	"chain_proto/peer"
	"chain_proto/transaction"
	"io/ioutil"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Repository struct {
	Block interface {
		GetByHash(hash [32]byte) (*block.Block, error)
		GetByHeight(height uint32) (*block.Block, error)
		GetByRange(offset uint32, limit uint32) ([]*block.Block, error)
		GetLatest() (*block.Block, error)
		Insert(b *block.Block) error
	}
	Tx interface {
		GetByBlockHash(blockHash [32]byte) ([]*transaction.Transaction, error)
		GetByHash(hash [32]byte) (*transaction.Transaction, error)
		BulkInsert(txs []*transaction.Transaction) error
	}
	Account interface {
		BulkInsert(accounts []*account.Account) error
		Insert(account *account.Account) error
		Get(addr string) (*account.Account, error)
	}
	Peer interface {
		GetAll() ([]*peer.Peer, error)
		Get(host string) (*peer.Peer, error)
		AddOrReplace(peer *peer.Peer) error
	}
}

func New() (*Repository, error) {
	dbPath := config.Config.DbPath
	dbDriver := config.Config.Driver
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
