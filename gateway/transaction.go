package gateway

import (
	gw "chain_proto/gateway/gw"
	"context"
	"errors"

	"github.com/golang/protobuf/ptypes/empty"
)

func (g *Gateway) GetTxsByBlockHash(_ context.Context, in *gw.GetTxByBlockHashRequest) (*gw.GetTransactionsResponse, error) {
	txs, err := g.bc.GetTxsByBlockHash(in.GetBlockHash())
	if err != nil {
		return nil, err
	}

	pbTxs := make([]*gw.Transaction, 0, len(txs))
	for _, tx := range txs {
		pbTxs = append(pbTxs, toPbTransaction(tx))
	}

	return &gw.GetTransactionsResponse{
		Transactions: pbTxs,
	}, nil
}

func (g *Gateway) GetTransactionByHash(_ context.Context, in *gw.GetTransactionByHashRequest) (*gw.GetTransactionResponse, error) {
	tx, err := g.bc.GetTransactionByHash(in.GetTxHash())
	if err != nil {
		return nil, err
	}

	return &gw.GetTransactionResponse{
		Transaction: toPbTransaction(tx),
	}, nil
}

func (g *Gateway) SendTransaction(_ context.Context, in *gw.SendTransactionRequest) (*empty.Empty, error) {
	tx, err := toTransaction(in.GetTransaction())
	if err != nil {
		return nil, err
	}
	if err := g.bc.AddTransaction(tx); err != nil {
		return nil, errors.New("error: Faliled to add a new transacation")
	}

	return &empty.Empty{}, nil
}
