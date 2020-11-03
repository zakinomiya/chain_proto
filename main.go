package main

import (
	"go_chain/common"
	"go_chain/config"
	"go_chain/server"
)

func main() {
	common.LogginSettings(config.Config.LogFile, config.Config.DefaultLogLevel)

	server := server.New(&config.Config)
	server.Start()

	// if err := server.Start(); err != nil {
	// 	log.Printf("Failed to start server. %s", err.Error())
	// 	os.Exit(1)
	// }

	// c := make(chan int)

	// for {
	// 	select {
	// 	case <-c:

	// 	}
	// }

}
