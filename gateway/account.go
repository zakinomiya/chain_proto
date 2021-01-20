package gateway

import (
	gw "chain_proto/gateway/gw"
	"context"
)

func (bs *BlockchainService) GetAccount(_ context.Context, in *gw.GetAccountRequest) (*gw.GetAccountResponse, error) {
	acc, err := bs.bc.GetAccount(in.GetAddr())
	if err != nil {
		return nil, err
	}

	return &gw.GetAccountResponse{
		Account: &gw.Account{
			Addr:    acc.Addr,
			Balance: acc.BalanceString()},
	}, nil
}
