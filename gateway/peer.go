package gateway

import (
	gw "chain_proto/gateway/gw"
	"context"

	"github.com/golang/protobuf/ptypes/empty"
)

func (g *Gateway) PropagtePeer(ctx context.Context, in *gw.SendPeerRequest) (*empty.Empty, error) {
	peer := in.GetPeer()

	if err := g.bc.AddPeer(peer.Host); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}
