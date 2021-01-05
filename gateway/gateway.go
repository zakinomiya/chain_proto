package gateway

import (
	"chain_proto/account"
	"chain_proto/block"
	"chain_proto/config"
	"chain_proto/transaction"
	"context"
	"log"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	gw "chain_proto/gateway/gw"
)

type Blockchain interface {
	GetBlockByHash(hash string) (*block.Block, error)
	GetBlockByHeight(height int32) (*block.Block, error)
	GetBlocks(offset int32, limit int32) ([]*block.Block, error)
	AddBlock(block *block.Block) bool
	GetLatestBlock() (*block.Block, error)
	GetTxsByBlockHash(blockHash string) ([]*transaction.Transaction, error)
	GetTransactionByHash(hash string) ([]*transaction.Transaction, error)
	AddTransaction(tx *transaction.Transaction) bool
	GetAccount(addr string) (*account.Account, error)
	AddAccount(account *account.Account) bool
	AddPeer(host string) bool
	Sync(offset int) ([]*block.Block, error)
}

type Gateway struct {
	httpServer *HTTPServer
	bc         Blockchain
}

func New(bc Blockchain) *Gateway {
	return &Gateway{
		httpServer: NewHTTPServer(config.Config.HTTP.Port),
		bc:         bc,
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

	if config.Config.RPC.Enabled {
		log.Println("info: starting rpc server")
	}

	if config.Config.HTTP.Enabled {
		log.Println("info: starting http server")
	}

	if config.Config.Websocket.Enabled {
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
