package gateway

import (
	gw "chain_proto/gateway/gw"
	"context"
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
