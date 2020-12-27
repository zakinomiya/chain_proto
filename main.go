package main

import (
	"go_chain/db/models"
	"io/ioutil"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, _ := sqlx.Open("sqlite3", "./data/blockchain.db")
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	sql, _ := ioutil.ReadFile("db/sql/insert_account.sql")
	log.Println(string(sql))

	// stmt, err := db.PrepareNamed(string(sql))
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// log.Println(stmt.QueryString)

	m := []*models.AccountModel{
		{Addr: "sss", Balance: 0},
		{Addr: "aaa", Balance: 10},
	}

	tx, _ := db.Beginx()

	log.Println(string(sql))
	if _, err := tx.NamedExec(string(sql), m); err != nil {
		log.Println(err)
		return
	}

	tx.Commit()
}
