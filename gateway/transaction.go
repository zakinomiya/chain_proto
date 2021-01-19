package gateway

import (
	"chain_proto/common"
	gw "chain_proto/gateway/gw"
	"context"
	"encoding/hex"
	"errors"

	"github.com/golang/protobuf/ptypes/empty"
)

func (bs *BlockchainService) GetTxsByBlockHash(_ context.Context, in *gw.GetTxsByBlockHashRequest) (*gw.GetTxsByBlockHashResponse, error) {
	blockHash, err := hex.DecodeString(in.GetBlockHash())
	if err != nil {
		return nil, err
	}

	txs, err := bs.bc.GetTxsByBlockHash(common.ReadByteInto32(blockHash))
	if err != nil {
		return nil, err
	}

	pbTxs := make([]*gw.Transaction, 0, len(txs))
	for _, tx := range txs {
		pbTxs = append(pbTxs, toPbTransaction(tx))
	}

	return &gw.GetTxsByBlockHashResponse{
		Transactions: pbTxs,
	}, nil
}

func (bs *BlockchainService) GetTransactionByHash(_ context.Context, in *gw.GetTransactionByHashRequest) (*gw.GetTransactionResponse, error) {
	hash, err := hex.DecodeString(in.GetTxHash())
	if err != nil {
		return nil, err
	}

	tx, err := bs.bc.GetTransactionByHash(common.ReadByteInto32(hash))
	if err != nil {
		return nil, err
	}

	return &gw.GetTransactionResponse{
		Transaction: toPbTransaction(tx),
	}, nil
}

func (bs *BlockchainService) PropagateTransaction(_ context.Context, in *gw.PropagateTransactionRequest) (*empty.Empty, error) {
	tx, err := toTransaction(in.GetTransaction())
	if err != nil {
		return nil, err
	}
	if err := bs.bc.AddTransaction(tx); err != nil {
		return nil, errors.New("error: Faliled to add a new transacation")
	}

	return &empty.Empty{}, nil
}
