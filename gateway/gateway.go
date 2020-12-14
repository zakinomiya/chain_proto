package gateway

import (
	"context"
	"go_chain/config"
	"log"
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
	if g.config.RPC.Enabled {
		log.Println("info: starting rpc server")
	}

	if g.config.HTTP.Enabled {
		log.Println("info: starting http server")
		go g.httpServer.Start()
	}

	if g.config.Websocket.Enabled {
		log.Println("info: starting rpc server")
	}

	log.Println("info: successfully started servers")
	return nil
}

func (g *Gateway) Stop() {
	if err := g.httpServer.Shutdown(context.Background()); err != nil {
		log.Fatalln("Failed to stop http server")
	}
}
