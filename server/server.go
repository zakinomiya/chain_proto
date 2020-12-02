package server

import (
	"errors"
	"go_chain/blockchain"
	"go_chain/config"
	"go_chain/miner"
	"go_chain/repository"
	"go_chain/wallet"
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

func New(config *config.ConfigSettings) (*Server, error) {
	if config.Settings.Miner.SecretKeyStr == "" {
		return nil, errors.New("No miner key provided")
	}

	repository := repository.New(config.Settings.DbPath, config.Settings.Driver)
	blockchain := blockchain.New(config.Settings.ChainID, repository)

	wal, err := wallet.RestoreWallet(config.Settings.Miner.SecretKeyStr)
	if err != nil {
		return nil, err
	}

	miner := miner.New(blockchain, wal)
	// gateway := &gateway.Gateway{}

	return &Server{config: config, services: []Service{repository, blockchain, miner}}, nil
}

func (server *Server) Start() error {
	log.Printf("info: Starting MVB node. ChainID=%d\n", server.config.Settings.ChainID)

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
