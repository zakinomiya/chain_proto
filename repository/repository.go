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

func (d *database) exec(filename string, data interface{}) error {
	args := data
	if args == nil {
		args = ""
	}

	stmt, err := readSQL(filename)
	if err != nil {
		return err
	}

	if _, err := d.db.NamedExec(stmt, data); err != nil {
		log.Printf("error: Failed to execute SQL. sql=%s err: %v\n", filename, err)
		return err
	}

	return nil
}

func (d *database) find(filename string, data interface{}) (*sqlx.Rows, error) {
	args := data
	if args == nil {
		args = ""
	}

	stmt, err := readSQL(filename)
	if err != nil {
		log.Printf("sql file not found. filename=%s\n", filename)
		return nil, err
	}

	rows, err := d.db.NamedQuery(stmt, args)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

type Repository struct {
	*database
}

func New(dbPath string, dbDriver string) *Repository {
	db := &database{
		dbPath: dbPath, dbDriver: dbDriver,
	}
	rep := &Repository{database: db}

	return rep
}

func (r *Repository) Start() error {
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

func (r *Repository) Stop() {
	r.db.Close()
}

func (r *Repository) ServiceName() string {
	return "Repository"
}

func readSQL(filename string) (string, error) {
	sql, err := ioutil.ReadFile("./repository/sql/" + filename)
	if err != nil {
		log.Printf("error: Error retrieving sql file. filename=%s\n", filename)
		return "", err
	}

	return string(sql), nil
}
