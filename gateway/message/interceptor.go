package message

import (
	"chain_proto/config"
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var MessagesCache map[string]bool

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if !strings.Contains(info.FullMethod, "Propagate") || !strings.Contains(info.FullMethod, "Sync") {
			return handler(ctx, req)
		}

		var messageID, chainID string

		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			chainID = md.Get("X-chain-id")[0]
			messageID = md.Get("X-message-id")[0]
		}

		if ok := checkChainID(chainID); !ok {
			return nil, status.Errorf(codes.InvalidArgument, "invalid chainID. ChainID is %s, but given %s.\n", config.Config.ChainID, chainID)
		}

		if checkMessageID(messageID) {
			return codes.OK, nil
		} else {
			MessagesCache[messageID] = true
		}

		ctx = context.WithValue(ctx, "MessageID", messageID)
		return handler(ctx, req)
	}
}

func checkChainID(chainID string) bool {
	return config.Config.ChainID == chainID
}

func checkMessageID(messageID string) bool {
	return MessagesCache[messageID]
}
