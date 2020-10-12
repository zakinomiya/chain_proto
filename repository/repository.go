package repository

import (
	"database/sql"
	"go_chain/block"
	"io/ioutil"
	"log"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

var repository *Repository
var once sync.Once

type Repository struct {
	dbPath          string
	dbDriver        string
	BlockRepository IBlockRepostitory
}

type IBlockRepostitory interface {
	Get(hash string) (*block.Block, error)
	GetAll() ([]*block.Block, error)
}

// func New(dbPath string, dbDriver string) *Repository {
func New() *Repository {
	once.Do(func() {
		repository = &Repository{dbPath: "./blockchain.db", dbDriver: "sqlite3", BlockRepository: newBlockRepository()}
	})

	return repository
}

func (r *Repository) Start() error {
	db, err := sql.Open(r.dbDriver, r.dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	initSQL, err := ioutil.ReadFile("./repository/sql/initialize.sql")
	if err != nil {
		log.Println("error: Error retrieving db initialise sql")
		return err
	}

	if _, err := db.Exec(string(initSQL)); err != nil {
		log.Printf("error: Failed to execute SQL. sql=initialize.sql err: %v\n", err)
		return err
	}

	return nil
}

func (r *Repository) Stop() {

}

func (r *Repository) ServiceName() string {
	return "Repository"
}
