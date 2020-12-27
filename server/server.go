package server

import (
	"errors"
	"go_chain/blockchain"
	"go_chain/config"
	"go_chain/db/repository"
	"go_chain/gateway"
	"go_chain/miner"
	"go_chain/wallet"
	"log"
)

type Service interface {
	Start() error
	Stop()
	ServiceName() string
}

type Server struct {
	config   *config.Configurations
	services []Service
}

func New(config *config.Configurations) (*Server, error) {
	if config.Miner.SecretKeyStr == "" {
		return nil, errors.New("No miner key provided")
	}

	repository, err := repository.New(config.DbPath, config.Driver)
	if err != nil {
		return nil, err
	}

	blockchain := blockchain.New(config.ChainID, repository)
	wal, err := wallet.RestoreWallet(config.Miner.SecretKeyStr)
	if err != nil {
		return nil, err
	}

	miner := miner.New(blockchain, wal, config.Enabled, config.Concurrent, config.MaxWorkersNum)
	gateway := gateway.New(&config.Network)

	return &Server{
		config:   config,
		services: []Service{blockchain, miner, gateway},
	}, nil
}

func (server *Server) Start() error {
	log.Printf("info: Starting MVB node. ChainID=%d\n", server.config.ChainID)

	for _, s := range server.services {
		stream := make(chan struct{})
		go func() {
			defer close(stream)
			if err := s.Start(); err != nil {
				return
			}
		}()
		// Wait until service starts
		select {
		case <-stream:
		}
		log.Printf("Successfully started service %s \n", s.ServiceName())
	}

	log.Println("Successfully started the node")
	return nil
}

func (server *Server) Stop() {
	log.Printf("info: Stopping MVB node. ChainID=%d\n", server.config.ChainID)

	for _, s := range server.services {
		log.Printf("Stopping service %s ...\n", s.ServiceName())
		s.Stop()
		log.Printf("Successfully stopped service %s \n", s.ServiceName())
	}

	log.Println("Successfully stopped the node")
}
