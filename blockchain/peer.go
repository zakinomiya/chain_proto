package blockchain

import (
	"chain_proto/block"
	"chain_proto/db/repository"
	"chain_proto/peer"
	"context"
	"io"
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

func (bc *Blockchain) Sync() error {
	ch := make(chan []*block.Block)
	errCh := make(chan error)

	go func() {
		if err := bc.c.Sync(context.Background(), ch, bc.CurrentBlockHeight()); err != nil {
			if err == io.EOF {
				return
			}

			errCh <- err
		}
	}()

	for {
		select {
		case blks := <-ch:
			for _, b := range blks {
				if err := bc.AddBlock(b); err != nil {
					return err
				}
			}
		case err := <-errCh:
			return err
		}
	}
}
