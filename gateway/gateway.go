package gateway

import (
	"context"
	"go_chain/config"
	"log"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	gw "go_chain/gateway/gw"
)

type Gateway struct {
	config     *config.Network
	httpServer *HTTPServer
}

func New(config *config.Network) *Gateway {
	return &Gateway{
		config:     config,
		httpServer: NewHTTPServer(config.HTTP.Port),
	}
}

func (g *Gateway) ServiceName() string {
	return "Gateway"
}

// Start starts servers.
// TODO run servers as goroutines
func (g *Gateway) Start() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	if g.config.RPC.Enabled {
		log.Println("info: starting rpc server")
	}

	if g.config.HTTP.Enabled {
		log.Println("info: starting http server")
	}

	if g.config.Websocket.Enabled {
		log.Println("info: starting rpc server")
	}

	log.Println("info: successfully started servers")

	go func() {
		if err := gw.RegisterHTTPServiceHandlerFromEndpoint(ctx, mux, "localhost:8081", opts); err != nil {
			log.Fatalln("Failed to start the http server")
		}
	}()

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	return nil
}

func (g *Gateway) Stop() {
	if err := g.httpServer.Shutdown(context.Background()); err != nil {
		log.Fatalln("Failed to stop http server")
	}
}
