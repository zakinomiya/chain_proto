package gateway

import (
	"chain_proto/common"
	"chain_proto/db/repository"
	gw "chain_proto/gateway/gw"
	"context"
	"encoding/hex"
	"log"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (bs *BlockchainService) GetLatestBlock(_ context.Context, in *gw.GetBlockByHashRequest) (*gw.GetBlockResponse, error) {
	blk := bs.bc.LatestBlock()
	pbBlk, err := toPbBlock(blk)
	if err != nil {
		return nil, err
	}

	return &gw.GetBlockResponse{
		Block: pbBlk,
	}, nil
}

func (bs *BlockchainService) GetBlockByHash(_ context.Context, in *gw.GetBlockByHashRequest) (*gw.GetBlockResponse, error) {
	blockHash, err := hex.DecodeString(in.GetBlockHash())
	if err != nil {
		log.Println("error: 3")
		return nil, err
	}

	blk, err := bs.bc.GetBlockByHash(common.ReadByteInto32(blockHash))
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, status.Errorf(codes.NotFound, "No block with block hash(%s) found", in.GetBlockHash())
		}

		return nil, err
	}

	pbBlk, err := toPbBlock(blk)
	if err != nil {
		log.Println("error: 1")
		return nil, err
	}

	return &gw.GetBlockResponse{
		Block: pbBlk,
	}, nil
}

func (bs *BlockchainService) GetBlockByHeight(_ context.Context, in *gw.GetBlockByHeightRequest) (*gw.GetBlockResponse, error) {
	blk, err := bs.bc.GetBlockByHeight(in.GetBlockHeight())
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

func (bs *BlockchainService) GetBlocks(_ context.Context, in *gw.GetBlocksRequest) (*gw.GetBlocksResponse, error) {
	blks, err := bs.bc.GetBlocks(in.GetOffset(), in.GetLimit())
	if err != nil {
		return nil, err
	}

	pbBlocks := make([]*gw.Block, 0)
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

func (bs *BlockchainService) PropagateBlock(_ context.Context, in *gw.PropagateBlockRequest) (*empty.Empty, error) {
	blk, err := toBlock(in.GetBlock())
	if err != nil {
		return nil, err
	}

	if err := bs.bc.AddBlock(blk); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}
