package gateway

import (
	"chain_proto/config"
	gw "chain_proto/gateway/gw"
	"context"
	"errors"
)

func (g *Gateway) GetBlockByHash(_ context.Context, in *gw.GetBlockByHashRequest) (*gw.Block, error) {
	b, err := g.bc.GetBlockByHash(in.GetBlockHash())
	if err != nil {
		return nil, err
	}

	return toPbBlock(b)
}

func (g *Gateway) GetBlockByHeight(_ context.Context, in *gw.GetBlockByHeightRequest) (*gw.Block, error) {
	b, err := g.bc.GetBlockByHeight(in.GetBlockHeight())
	if err != nil {
		return nil, err
	}
	return toPbBlock(b)
}

func (g *Gateway) GetBlocks(_ context.Context, in *gw.GetBlocksRequest) ([]*gw.Block, error) {
	blks, err := g.bc.GetBlocks(in.GetOffset(), in.GetLimit())
	if err != nil {
		return nil, err
	}

	pbBlocks := make([]*gw.Block, 0, len(blks))
	for _, b := range blks {
		pbBlock, err := toPbBlock(b)
		if err != nil {
			return nil, err
		}
		pbBlocks = append(pbBlocks, pbBlock)
	}

	return pbBlocks, nil
}

func (g *Gateway) AddBlock(_ context.Context, in *gw.SendBlockRequest) error {
	if in.GetMetadata().ChainID != config.ChainInfo.ChainID {
		return errors.New("error: Invalid chain Id.")
	}

	blk, err := toBlock(in.GetBlock())
	if err != nil {
		return err
	}

	if !g.bc.AddBlock(blk) {
		return err
	}

	return nil
}
