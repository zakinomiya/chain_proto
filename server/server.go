package server

import (
	"fmt"
	"go_chain/block"
	"go_chain/blockchain"
	"go_chain/config"
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
	firstBlock.SetAmount(100)
	firstBlock.SetHash("First Block")
	fmt.Printf("New block. %#v", firstBlock)

	fmt.Println("Adding new block")
	server.blockchain.AddNewBlock(firstBlock)

	fmt.Printf("Blocks in the blockchain: %#v", server.blockchain.Blocks())
}
