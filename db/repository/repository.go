package repository

import (
	"io/ioutil"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type database struct {
	db       *sqlx.DB
	dbPath   string
	dbDriver string
}

type Repository struct {
	*database
}

func New(dbPath string, dbDriver string) *Repository {
	return &Repository{
		database: &database{
			dbPath:   dbPath,
			dbDriver: dbDriver,
		},
	}
}

func (r *Repository) Connect() error {
	log.Printf("debug: DBInfo driver=%s, dbpath=%s\n", r.dbDriver, r.dbPath)
	db, err := sqlx.Open(r.dbDriver, r.dbPath)
	if err != nil {
		log.Println(err)
		return err
	}
	r.db = db

	err = r.exec("initialize.sql", map[string]interface{}{})
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) Close() {
	r.db.Close()
}

func readSQL(filename string) (string, error) {
	sql, err := ioutil.ReadFile("./repository/sql/" + filename)
	if err != nil {
		log.Printf("error: Error retrieving sql file. filename=%s\n", filename)
		return "", err
	}

	return string(sql), nil
}
