package repository

import (
	"io/ioutil"
	"log"

	"github.com/jmoiron/sqlx"
)

type database struct {
	db       *sqlx.DB
	dbPath   string
	dbDriver string
}

type prepareFunc func(query string) (*sqlx.NamedStmt, error)

func (d *database) disconnect() {
	d.db.Close()
}

// prepareNamed reads a sql file and preapare named query with it
func (d *database) prepareNamed(fn prepareFunc, filename string) (*sqlx.NamedStmt, error) {
	sql, err := ioutil.ReadFile("./db/sql/" + filename)
	if err != nil {
		log.Printf("error: Error retrieving sql file. filename=%s\n", filename)
		return nil, err
	}

	return fn(string(sql))
}

// queryRow queries a single row
func (d *database) queryRow(filename string, data interface{}) (*sqlx.Row, error) {
	args := data
	if args == nil {
		args = map[string]interface{}{}
	}

	stmt, err := d.prepareNamed(d.db.PrepareNamed, filename)
	if err != nil {
		log.Printf("sql file not found. filename=%s\n", filename)
		return nil, err
	}

	return stmt.QueryRowx(args), nil
}

// queryRows queries multiple rows
func (d *database) queryRows(filename string, data interface{}) (*sqlx.Rows, error) {
	args := data
	if args == nil {
		args = map[string]interface{}{}
	}

	stmt, err := d.prepareNamed(d.db.PrepareNamed, filename)
	if err != nil {
		log.Printf("error: failed to prepare sql statement. err=%v\n", err)
		return nil, err
	}

	rows, err := stmt.Queryx(args)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// command executes sql statement except select
func (d *database) command(filename string, data interface{}) error {
	args := data
	if args == nil {
		args = ""
	}

	stmt, err := d.prepareNamed(d.db.PrepareNamed, filename)
	if err != nil {
		return err
	}

	if _, err := stmt.Exec(args); err != nil {
		log.Printf("error: Failed to execute SQL. sql=%s err: %v\n", filename, err)
		return err
	}

	return nil
}

// transaction receives several actions(sql statement) and exectues them as a transaction.
func (d *database) txCommand(tx *sqlx.Tx, filename string, data interface{}) error {
	stmt, err := d.prepareNamed(tx.PrepareNamed, filename)
	if err != nil {
		return err
	}

	sql, _ := ioutil.ReadFile("./db/sql/" + filename)

	log.Printf("debug: action=txCommand. sql=%s", string(sql))
	if _, err := tx.NamedExec(string(sql), data); err != nil {
		log.Printf("debug: action=txCommand. sql=%s", stmt.QueryString)
		return err
	}

	log.Println("debug: action=txCommand. filename=", filename)
	return nil
}
