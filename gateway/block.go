package gateway

import (
	gw "chain_proto/gateway/gw"
	"context"

	"github.com/golang/protobuf/ptypes/empty"
)

func (g *Gateway) GetLatestBlock(_ context.Context, in *gw.GetBlockByHashRequest) (*gw.GetBlockResponse, error) {
	blk, err := g.bc.GetLatestBlock()
	if err != nil {
		return nil, err
	}

	pbBlk, err := toPbBlock(blk)
	if err != nil {
		return nil, err
	}

	return &gw.GetBlockResponse{
		Block: pbBlk,
	}, nil
}

func (g *Gateway) GetBlockByHash(_ context.Context, in *gw.GetBlockByHashRequest) (*gw.GetBlockResponse, error) {
	blk, err := g.bc.GetBlockByHash(in.GetBlockHash())
	if err != nil {
		return nil, err
	}

	pbBlk, err := toPbBlock(blk)
	if err != nil {
		return nil, err
	}

	return &gw.GetBlockResponse{
		Block: pbBlk,
	}, nil
}

func (g *Gateway) GetBlockByHeight(_ context.Context, in *gw.GetBlockByHeightRequest) (*gw.GetBlockResponse, error) {
	blk, err := g.bc.GetBlockByHeight(in.GetBlockHeight())
	if err != nil {
		return nil, err
	}

	pbBlk, err := toPbBlock(blk)
	if err != nil {
		return nil, err
	}

	return &gw.GetBlockResponse{
		Block: pbBlk,
	}, nil
}

func (g *Gateway) GetBlocks(_ context.Context, in *gw.GetBlocksRequest) (*gw.GetBlocksResponse, error) {
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

	return &gw.GetBlocksResponse{
		Blocks: pbBlocks,
	}, nil
}

func (g *Gateway) SendBlock(_ context.Context, in *gw.SendBlockRequest) (*empty.Empty, error) {
	blk, err := toBlock(in.GetBlock())
	if err != nil {
		return nil, err
	}

	if err := g.bc.AddBlock(blk); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}
