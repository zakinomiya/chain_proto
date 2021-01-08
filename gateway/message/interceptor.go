package message

import (
	"chain_proto/config"
	"context"
	"log"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if !strings.Contains(info.FullMethod, "Propagate") || !strings.Contains(info.FullMethod, "Sync") {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			log.Println("error: failed to retrieve information from header")
			return nil, status.Error(codes.Internal, "request failed")
		}

		chainID := md.Get("X-chain-id")[0]
		if chainID != config.Config.ChainID {
			return nil, status.Errorf(codes.InvalidArgument, "invalid chainID. ChainID is %s, but given %s.\n", config.Config.ChainID, chainID)
		}

		messageID := md.Get("X-message-id")[0]
		ctx = context.WithValue(ctx, "MessageID", messageID)

		return handler(ctx, req)
	}
}
