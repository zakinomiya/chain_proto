package peer

type Peer struct {
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
