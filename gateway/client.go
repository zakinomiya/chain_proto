package gateway

import (
	"chain_proto/block"
	gw "chain_proto/gateway/gw"
	"chain_proto/gateway/message"
	"chain_proto/peer"
	"chain_proto/transaction"
	"context"
	"time"

	"google.golang.org/grpc"
)

const (
	requestTimeout = 15 * time.Second
)

type Client struct {
	neighbours []*peer.Peer
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) tidyUp() {
	newNeighbours := make([]*peer.Peer, 0)
	for _, p := range c.neighbours {
		if p.FailCount <= 15 {
			newNeighbours = append(newNeighbours, p)
		}
	}
}

func (c *Client) AddNeighbour(p *peer.Peer) {
	c.neighbours = append(c.neighbours, p)
}

func (c *Client) PropagateBlock(b *block.Block) error {
	defer c.tidyUp()
	pbBlock, err := toPbBlock(b)
	if err != nil {
		return err
	}

	invoke := func(p *peer.Peer) {
		ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
		defer cancel()

		conn, err := p.Connect(grpc.WithInsecure(), grpc.WithUnaryInterceptor(message.UnaryClientInterceptor()))
		if err != nil {
			return
		}

		s := gw.NewBlockchainServiceClient(conn)
		if _, err := s.PropagateBlock(ctx, &gw.PropagateBlockRequest{Block: pbBlock}, nil); err != nil {
			return
		}
	}

	for _, p := range c.neighbours {
		go invoke(p)
	}

	return nil
}

func (c *Client) PropagateTransaction(tx *transaction.Transaction) error {
	pbTransaction := toPbTransaction(tx)

	req := func(p *peer.Peer) {
		ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
		defer cancel()

		conn, err := p.Connect(grpc.WithInsecure(), grpc.WithUnaryInterceptor(message.UnaryClientInterceptor()))
		if err != nil {
			return
		}

		s := gw.NewBlockchainServiceClient(conn)
		if _, err := s.PropagateTransaction(ctx, &gw.PropagateTransactionRequest{Transaction: pbTransaction}, nil); err != nil {
			return
		}
	}

	c.invoke(req)
	return nil
}

func (c *Client) invoke(reqFunc func(p *peer.Peer)) {
	defer c.tidyUp()
	for _, p := range c.neighbours {
		go reqFunc(p)
	}
}
