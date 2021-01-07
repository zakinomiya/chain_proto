package gateway

import (
	gw "chain_proto/gateway/gw"
	"context"

	"google.golang.org/grpc/metadata"
)

func (g *Gateway) GetAccount(_ context.Context, in *gw.GetAccountRequest) (*gw.GetAccountResponse, error) {
	acc, err := g.bc.GetAccount(in.GetAddr())
	if err != nil {
		return nil, err
	}

	return &gw.GetAccountResponse{
		Account: &gw.Account{
			Addr:    acc.Addr,
			Balance: acc.Balance,
		},
	}, nil
}

func (g *Gateway) SendAccount(ctx context.Context, in *gw.GetAccountRequest) (*gw.GetAccountResponse, error) {
	// Check metadata

	meta, err := metadata.FromIncomingContext(ctx)

}
