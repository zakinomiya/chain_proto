package main

import (
	"fmt"
	"go_chain/blockchain"
	"go_chain/config"
	"go_chain/miner"
	"go_chain/utils"
)

func main() {
	utils.LogginSettings(config.Config.LogFile, config.Config.DefaultLogLevel)

	// server := server.New(&config.Config)

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

	miner := miner.New(&blockchain.Blockchain{})
	b := miner.CalcGenesis()
	fmt.Println(fmt.Sprintf("%x", b.Hash()), b.Nonce(), fmt.Sprintf("%x", b.MerkleRoot()), b.Timestamp())
}
