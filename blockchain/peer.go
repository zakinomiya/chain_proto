package blockchain

import (
	"chain_proto/db/repository"
	"chain_proto/peer"
)

// AddPeer registers a new peer to the client
func (bc *Blockchain) AddPeer(addr string, network string) error {
	p := peer.New(addr, network)
	bc.c.AddNeighbour(p)
	return bc.r.Peer.AddOrReplace(p)
}

// GetPeers returns the registered peers.
// If no peers are found, an empty slice will be returned.
func (bc *Blockchain) GetPeers() ([]*peer.Peer, error) {
	peers, err := bc.r.Peer.GetAll()
	if err != nil {
		if err == repository.ErrNotFound {
			return make([]*peer.Peer, 0), nil
		}
		return nil, err
	}

	return peers, nil
}
