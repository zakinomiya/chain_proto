package gateway

import (
	gw "chain_proto/gateway/gw"
	"context"
	"errors"
)

func (g *Gateway) GetTxsByBlockHash(_ context.Context, in *gw.GetTxByBlockHashRequest) ([]*gw.Transaction, error) {
	txs, err := g.bc.GetTxsByBlockHash(in.GetBlockHash())
	if err != nil {
		return nil, err
	}

	pbTxs := make([]*gw.Transaction, 0, len(txs))
	for _, tx := range txs {
		pbTxs = append(pbTxs, toPbTransaction(tx))
	}

	return pbTxs, nil
}

func (g *Gateway) GetTransactionByHash(_ context.Context, in *gw.GetTransactionByHashRequest) (*gw.Transaction, error) {
	tx, err := g.bc.GetTransactionByHash(in.GetTxHash())
	if err != nil {
		return nil, err
	}

	return toPbTransaction(tx), nil
}

func (g *Gateway) SendTransaction(_ context.Context, in *gw.SendTransactionRequest) error {
	tx, err := toTransaction(in.GetTransaction())
	if err != nil {
		return err
	}
	if !g.bc.AddTransaction(tx) {
		return errors.New("error: Faliled to add a new transacation")
	}

	return nil
}
