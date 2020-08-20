package main

import (
	"go_chain/config"
	"go_chain/server"
	"go_chain/utils"
	"log"
	"os"
)

func main() {
	utils.LogginSettings(config.Config.LogFile)

	server := server.New(&config.Config)

	if err := server.Start(); err != nil {
		log.Printf("Failed to start server. %s", err.Error())
		os.Exit(1)
	}
}
