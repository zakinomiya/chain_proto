package peer

import (
	"time"

	"google.golang.org/grpc"
)

const (
	requestTimeout = 15 * time.Second
)

type Peer struct {
	conn      *grpc.ClientConn
	addr      string
	network   string
	FailCount int
}

func New(addr string, network string) *Peer {
	return &Peer{
		addr:      addr,
		network:   network,
		FailCount: 0,
	}
}

func (p *Peer) Addr() string {
	return p.addr
}

func (p *Peer) Network() string {
	return p.network
}

func (p *Peer) Connect(opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	if p.conn == nil {
		conn, err := grpc.Dial(p.Addr(), opts...)
		if err != nil {
			p.FailCount += 1
			return nil, err
		}
		p.conn = conn
	}
	return p.conn, nil
}
