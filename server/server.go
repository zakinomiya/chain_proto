package server

import (
	"go_chain/blockchain"
	"go_chain/config"
	"go_chain/gateway"
	"go_chain/miner"
	"go_chain/repository"
	"log"
)

type Service interface {
	Start() error
	Stop()
	ServiceName() string
}

type Server struct {
	config   *config.ConfigSettings
	services []Service
}

func New(config *config.ConfigSettings) *Server {
	r := repository.New(config.Path, config.Driver)
	bc := blockchain.New(config.ChainID, r)
	m := miner.New(bc)
	g := &gateway.Gateway{}

	return &Server{config: config, services: []Service{r, bc, m, g}}
}

func (server *Server) Start() error {
	log.Printf("info: Starting MVB node. ChainID=%d\n", server.config.ChainID)

	for _, s := range server.services {
		log.Printf("Starting service %s \n", s.ServiceName())
		if err := s.Start(); err != nil {
			log.Printf("Failed to start service %s \n", s.ServiceName())
			return err
		}
		log.Printf("Successfully started service %s \n", s.ServiceName())
	}

	log.Println("Successfully started the node")

	return nil
}
