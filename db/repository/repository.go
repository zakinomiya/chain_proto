package repository

import (
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
		GetTxByBlockHash(blockHash [32]byte) ([]*transaction.Transaction, error)
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

	err = database.command("initialize.sql", map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	return &Repository{
		Block: &BlockRepository{database: database},
	}, nil
}

func readSQL(filename string) (string, error) {
	sql, err := ioutil.ReadFile("./repository/sql/" + filename)
	if err != nil {
		log.Printf("error: Error retrieving sql file. filename=%s\n", filename)
		return "", err
	}

	return string(sql), nil
}
