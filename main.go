package main

import (
	"go_chain/block"
	"go_chain/blockchain"
	"go_chain/config"
	"go_chain/utils"
	"log"
)

func main() {
	utils.LogginSettings(config.Config.LogFile)

	blockchain := blockchain.New(nil)

	log.Printf("Initialised Blockchain: %#v \n", blockchain.Blocks())

	firstBlock := block.New()
	firstBlock.SetAmount(200)
	firstBlock.SetHash("First Block")

	blockchain.AddNewBlock(firstBlock)

	log.Printf("Initialised Blockchain: %#v \n", blockchain.Blocks())
}
