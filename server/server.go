package server

import (
	"chain_proto/blockchain"
	"chain_proto/config"
	"chain_proto/db/repository"
	"chain_proto/gateway"
	"chain_proto/miner"
	"chain_proto/peer"
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

	r, err := repository.New(config.Config.DbPath, config.Config.Driver, config.Config.SQLPath)
	if err != nil {
		return nil, err
	}

	seedNodes := []*peer.Peer{}
	for _, s := range config.Config.Network.Seeds {
		p := peer.New(s.Addr, s.Network)
		seedNodes = append(seedNodes, p)
	}

	bc := blockchain.New(r, gateway.NewClient(seedNodes...))
	wal, err := wallet.RestoreWallet(config.Config.SecretKeyStr)
	if err != nil {
		return nil, err
	}

	miner := miner.New(bc, wal)
	gateway := gateway.New(bc)
	bc.SetMiner(miner)

	return &Server{
		services: []Service{bc, miner, gateway},
	}, nil
}

func (server *Server) Start() error {
	log.Printf("info: Starting MVB node. ChainID=%s\n", config.Config.ChainID)

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
	log.Printf("info: Stopping MVB node. ChainID=%s\n", config.Config.ChainID)

	for _, s := range server.services {
		log.Printf("Stopping service %s ...\n", s.ServiceName())
		s.Stop()
		log.Printf("Successfully stopped service %s \n", s.ServiceName())
	}

	log.Println("Successfully stopped the node")
}
