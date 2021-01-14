package message

import (
	"chain_proto/config"
	"context"
	"log"
	"strings"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		log.Printf("debug: new request received. funcName=%s\n", info.FullMethod)
		if !strings.Contains(info.FullMethod, "Propagate") || !strings.Contains(info.FullMethod, "Sync") {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			log.Println("error: failed to retrieve information from header")
			return nil, status.Error(codes.Internal, "request failed")
		}

		chainID := md.Get("chainID")[0]
		if chainID != config.Config.ChainID {
			return nil, status.Errorf(codes.InvalidArgument, "invalid chainID. ChainID is %s, but given %s.\n", config.Config.ChainID, chainID)
		}

		messageID := md.Get("messageID")[0]
		ctx = context.WithValue(ctx, "MessageID", messageID)

		return handler(ctx, req)
	}
}

func UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		log.Printf("debug: new outgoing grpc request. endpoint=%s. method=%s\n", cc.Target(), method)

		ctx = metadata.AppendToOutgoingContext(ctx, "chainID", config.Config.ChainID, "messageID", genMessageID())
		if err := invoker(ctx, method, req, reply, cc, opts...); err != nil {
			log.Printf("debug: invoke method(%s) request to %s failed\n", method, cc.Target())
			return err
		}
		log.Printf("debug: invoke method(%s) request to %s succeeded\n", method, cc.Target())
		return nil
	}
}

func genMessageID() string {
	return uuid.New().String()
}
