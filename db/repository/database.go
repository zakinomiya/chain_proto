package repository

import (
	"log"

	"github.com/jmoiron/sqlx"
)

type database struct {
	db       *sqlx.DB
	dbPath   string
	dbDriver string
}

func (d *database) disconnect() {
	d.db.Close()
}

func (d *database) query(filename string, data interface{}) (*sqlx.Rows, error) {
	args := data
	if args == nil {
		args = map[string]interface{}{}
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

func (d *database) command(filename string, data interface{}) error {
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
