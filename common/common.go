package common

import (
	"io/ioutil"
	"log"
	"strconv"
)

func IntToByteSlice(b int) []byte {
	return []byte(strconv.Itoa(b))
}

func ReadSQL(filename string) (string, error) {
	sql, err := ioutil.ReadFile("./repository/sql/" + filename)
	if err != nil {
		log.Printf("error: Error retrieving sql file. filename=%s\n", filename)
		return "", err
	}

	return string(sql), nil
}
