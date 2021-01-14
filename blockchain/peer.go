package blockchain

import (
	"chain_proto/db/repository"
	"chain_proto/peer"
)

func (bc *Blockchain) AddPeer(addr string, network string) error {
	p := peer.New(addr, network)
	bc.client.AddNeighbour(p)
	return bc.repository.Peer.AddOrReplace(p)
}

func (bc *Blockchain) GetPeers() ([]*peer.Peer, error) {
	peers, err := bc.repository.Peer.GetAll()
	if err != nil {
		if err == repository.ErrNotFound {
			return make([]*peer.Peer, 0), nil
		}
		return nil, err
	}

	return peers, nil
}
