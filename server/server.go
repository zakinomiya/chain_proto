package server

import (
	"fmt"
	"go_chain/block"
	"go_chain/blockchain"
	"go_chain/config"
	"go_chain/utils"
	"log"
)

type Service interface {
	Start() error
	Stop()
	ServiceName() string
}

type Server struct {
	config     *config.ConfigSettings
	blockchain *blockchain.Blockchain
}

func New(config *config.ConfigSettings) *Server {
	blockchain := blockchain.New(nil)

	return &Server{config, blockchain}
}

func (server *Server) Start() error {

	services := []Service{
		server.blockchain,
	}

	for _, s := range services {
		log.Printf("Starting service %s \n", s.ServiceName())
		if err := s.Start(); err != nil {
			log.Printf("Failed to start service %s \n", s.ServiceName())
			return err
		}
		log.Printf("Successfully started service %s \n", s.ServiceName())
	}

	log.Println("Successfully started the node")

	server.test()
	return nil
}

func (server *Server) test() {
	firstBlock := block.New()
	fmt.Printf("New block. %#v", firstBlock)

	fmt.Println("Adding coinbase transaction")
	tx := utils.NewCoinbase("some pubkey", 250)
	fmt.Printf("Transaction hash: %x \n", tx.Hash())
	firstBlock.AddTransaction(tx)

	server.blockchain.Mining(firstBlock)
}
