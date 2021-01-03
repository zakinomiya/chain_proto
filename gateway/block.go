package gateway

import (
	gw "chain_proto/gateway/gw"
	"context"
)

func (g *Gateway) GetBlockByHash(ctx context.Context, req *gw.GetBlockByHashRequest) (*gw.GetBlockResponse, error) {
}
