package gateway

import (
	gw "chain_proto/gateway/gw"
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func (bs *BlockchainService) Connect(ctx context.Context, in *empty.Empty) (*gw.ConnectResponse, error) {
	client, ok := peer.FromContext(ctx)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "No Peer information found")
	}

	if err := bs.bc.AddPeer(client.Addr.String(), client.Addr.Network()); err != nil {
		return nil, err
	}

	peers, _ := bs.bc.GetPeers()

	neighbours := make([]*gw.Peer, 0)
	for _, p := range peers {
		if p.Addr() == client.Addr.String() {
			continue
		}
		neighbours = append(neighbours, toPbPeer(p))
	}

	return &gw.ConnectResponse{
		Neighbours: neighbours,
	}, nil
}

func (bs *BlockchainService) Sync(in *gw.SyncRequest, server gw.BlockchainService_SyncServer) error {

	return nil
}
