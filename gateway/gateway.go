package gateway

import (
	"log"
	"time"
)

type Gateway struct {
	httpServer *HTTPServer
}

func New() *Gateway {
	return &Gateway{
		httpServer: NewHTTPServer(),
	}

}

func (g *Gateway) ServiceName() string {
	return "Gateway"
}

/// Start starts servers.
/// TODO run servers as goroutines
func (g *Gateway) Start() error {
	log.Println("info: starting http server")
	go g.httpServer.Start("8000")

	time.Sleep(1 * time.Second)

	log.Println("info: successfully started servers")
	return nil
}

func (g *Gateway) Stop() {

}
