package repository

import (
	"chain_proto/account"
	"chain_proto/block"
	"chain_proto/peer"
	"chain_proto/transaction"
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
		Insert(tx *transaction.Transaction) error
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

func New(dbPath string, dbDriver string, sqlPath string) (*Repository, error) {
	log.Printf("debug: DBInfo driver=%s, dbpath=%s\n", dbDriver, dbPath)
	db, err := sqlx.Open(dbDriver, dbPath)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	d := &database{
		db:       db,
		dbPath:   dbPath,
		dbDriver: dbDriver,
		sqlPath:  sqlPath,
	}

	sql, err := d.rawSQL("initialize.sql")
	if err != nil {
		return nil, err
	}

	_, err = d.db.Exec(sql)
	if err != nil {
		return nil, err
	}

	ar := &AccountRepository{database: d}
	tr := &TxRepository{database: d}

	return &Repository{
		Block: &BlockRepository{
			database:    d,
			account:     ar,
			transcation: tr,
		},
		Tx:      tr,
		Account: ar,
	}, nil
}
