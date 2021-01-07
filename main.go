package main

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var m map[string]bool

	fmt.Println(m["test"])
}
