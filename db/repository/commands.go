package repository

import (
	"log"

	"github.com/jmoiron/sqlx"
)

type commands struct {
	database *database
}

func newCommands() *commands {
	return &commands{}
}

func (c *commands) insert() {

}

func (c *commands) bulkInsert() {

}

func (c *commands) update() {

}

func (c *commands) bulkUpdate() {

}

func (c *commands) delete() {

}

func (c *commands) find(filename string, data interface{}) (*sqlx.Rows, error) {
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

func (c *commands) findAll() {

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
