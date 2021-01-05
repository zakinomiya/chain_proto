package server

import (
	"chain_proto/blockchain"
	"chain_proto/config"
	"chain_proto/db/repository"
	"chain_proto/gateway"
	"chain_proto/miner"
	"chain_proto/wallet"
	"errors"
	"log"
)

type Service interface {
	Start() error
	Stop()
	ServiceName() string
}

type Server struct {
	services []Service
}

func New() (*Server, error) {
	if config.Config.SecretKeyStr == "" {
		return nil, errors.New("No miner key provided")
	}

	repository, err := repository.New()
	if err != nil {
		return nil, err
	}

	blockchain := blockchain.New(repository)
	wal, err := wallet.RestoreWallet(config.Config.SecretKeyStr)
	if err != nil {
		return nil, err
	}

	miner := miner.New(blockchain, wal)
	gateway := gateway.New(blockchain)

	return &Server{
		services: []Service{blockchain, miner, gateway},
	}, nil
}

func (server *Server) Start() error {
	log.Printf("info: Starting MVB node. ChainID=%d\n", config.Config.ChainID)

	for _, s := range server.services {
		log.Printf("info: Staring service(%s)\n", s.ServiceName())
		if err := s.Start(); err != nil {
			return err
		}
		log.Printf("Successfully started service %s \n", s.ServiceName())
	}

	log.Println("Successfully started the node")
	return nil
}

func (server *Server) Stop() {
	log.Printf("info: Stopping MVB node. ChainID=%d\n", config.Config.ChainID)

	for _, s := range server.services {
		log.Printf("Stopping service %s ...\n", s.ServiceName())
		s.Stop()
		log.Printf("Successfully stopped service %s \n", s.ServiceName())
	}

	log.Println("Successfully stopped the node")
}
