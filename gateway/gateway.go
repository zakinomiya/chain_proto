package gateway

import (
	"chain_proto/config"
	gw "chain_proto/gateway/gw"
	"chain_proto/gateway/message"
	"context"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Gateway struct {
	grpcPort   string
	httpPort   string
	grpcServer *grpc.Server
	httpServer *http.Server
}

func New(bc Blockchain) *Gateway {
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(message.UnaryServerInterceptor()),
	)

	blockchainService := NewBlockchainService(bc)
	gw.RegisterBlockchainServiceServer(grpcServer, blockchainService)
	reflection.Register(grpcServer)

	httpServer := &http.Server{}

	return &Gateway{
		httpPort:   ":" + config.Config.HTTP.Port,
		grpcPort:   ":" + config.Config.RPC.Port,
		grpcServer: grpcServer,
		httpServer: httpServer,
	}
}

func (g *Gateway) ServiceName() string {
	return "Gateway"
}

// Start starts servers.
// TODO run servers as goroutines
func (g *Gateway) Start() error {
	log.Println("info: Opening the gate")

	if err := g.startGrpcServer(); err != nil {
		return err
	}

	g.startHTTPServer()
	log.Println("info: Successfully opened the gate")
	return nil
}

func (g *Gateway) startGrpcServer() error {
	log.Println("info: Starting grpc server")
	lis, err := net.Listen("tcp", g.grpcPort)
	if err != nil {
		return err
	}

	go func() {
		if err := g.grpcServer.Serve(lis); err != nil {
			log.Fatalln("fatal: Failed to start grpc server. err=", err)
		}
	}()

	return nil
}

func (g *Gateway) startHTTPServer() error {
	log.Println("info: Starting http server")
	mux := runtime.NewServeMux()
	options := []grpc.DialOption{grpc.WithInsecure()}

	if err := gw.RegisterBlockchainServiceHandlerFromEndpoint(context.Background(), mux, g.grpcPort, options); err != nil {
		return err
	}

	g.httpServer.Addr = g.httpPort
	g.httpServer.Handler = mux

	go func() {
		if err := g.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalln("fatal: failed to start http server. err=", err)
		}
	}()

	return nil
}

func (g *Gateway) Stop() {
	log.Println("info: stopping gateway")
	g.httpServer.Shutdown(context.Background())
	g.grpcServer.GracefulStop()
	log.Println("info: successfully stopped all the server")
}
