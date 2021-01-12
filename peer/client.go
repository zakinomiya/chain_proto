package peer

import (
	"google.golang.org/grpc"
)

type Client struct {
	connections []*grpc.ClientConn
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Connect(peers []*Peer) error {
	for _, p := range peers {
		conn, err := grpc.Dial(p.Addr(), []grpc.DialOption{grpc.WithInsecure()}...)
		if err != nil {
			return err
		}
		c.connections = append(c.connections, conn)
	}
	return nil
}
