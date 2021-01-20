package gateway

import (
	"chain_proto/account"
	"chain_proto/block"
	gw "chain_proto/gateway/gw"
	"chain_proto/gateway/message"
	"chain_proto/peer"
	"chain_proto/transaction"
	"context"
	"errors"
	"fmt"
	"log"
	"sort"
	"sync"
	"time"

	"google.golang.org/grpc"
)

const (
	requestTimeout = 15 * time.Second
)

type withReturnFn func(peer *peer.Peer) (interface{}, error)
type withoutReturnFn func(peer *peer.Peer) error

type Client struct {
	m          *sync.Mutex
	neighbours []*peer.Peer
}

func NewClient() *Client {
	return &Client{
		m: &sync.Mutex{},
	}
}

func (c *Client) tidyUp() {
	c.m.Lock()
	defer c.m.Unlock()

	newNeighbours := make([]*peer.Peer, 0)
	for _, p := range c.neighbours {
		if p.FailCount <= 15 {
			newNeighbours = append(newNeighbours, p)
		}
	}

	c.neighbours = newNeighbours
}

func (c *Client) sortNeghbours() {
	c.m.Lock()
	defer c.m.Unlock()

	sort.Slice(c.neighbours, func(i, j int) bool { return c.neighbours[i].FailCount < c.neighbours[i].FailCount })
}

func (c *Client) AddNeighbour(p *peer.Peer) {
	c.neighbours = append(c.neighbours, p)
}

func (c *Client) PropagateBlock(ctx context.Context, b *block.Block) error {
	defer c.tidyUp()
	pbBlock, err := toPbBlock(b)
	if err != nil {
		return err
	}

	req := func(p *peer.Peer) error {
		conn, err := p.Connect(grpc.WithInsecure(), grpc.WithUnaryInterceptor(message.UnaryClientInterceptor()))
		if err != nil {
			return err
		}

		s := gw.NewBlockchainServiceClient(conn)
		if _, err := s.PropagateBlock(ctx, &gw.PropagateBlockRequest{Block: pbBlock}); err != nil {
			return err
		}

		return nil
	}

	c.propagate(req)
	return nil
}

func (c *Client) PropagateTransaction(tx *transaction.Transaction) error {
	pbTransaction := toPbTransaction(tx)

	req := func(p *peer.Peer) error {
		ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
		defer cancel()

		conn, err := p.Connect(grpc.WithInsecure(), grpc.WithUnaryInterceptor(message.UnaryClientInterceptor()))
		if err != nil {
			log.Println("error:", err)
			return err
		}

		s := gw.NewBlockchainServiceClient(conn)
		_, err = s.PropagateTransaction(ctx, &gw.PropagateTransactionRequest{Transaction: pbTransaction})
		if err != nil {
			return err
		}

		return nil
	}

	c.propagate(req)
	return nil
}

func (c *Client) GetBlockByHash(ctx context.Context, hash [32]byte) (*block.Block, error) {
	req := func(p *peer.Peer) (interface{}, error) {
		conn, err := p.Connect(grpc.WithInsecure(), grpc.WithUnaryInterceptor(message.UnaryClientInterceptor()))
		if err != nil {
			return nil, fmt.Errorf("err: failed to connect to peer(%s) %s\n", p.Addr(), err)
		}
		defer conn.Close()

		s := gw.NewBlockchainServiceClient(conn)
		res, err := s.GetBlockByHash(ctx, &gw.GetBlockByHashRequest{BlockHash: fmt.Sprintf("%x", hash)})
		if err != nil {
			return nil, fmt.Errorf("%+v. err: failed to get response. action=GetBlockByHash.\n", err)
		}

		return res, nil
	}

	res, err := c.invoke(req)
	if err != nil {
		return nil, err
	}

	blk, err := toBlock(res.(*gw.GetBlockResponse).GetBlock())
	if err != nil {
		return nil, err
	}

	return blk, nil
}

func (c *Client) GetBlockByHeight(ctx context.Context, height uint32) (*block.Block, error) {
	req := func(p *peer.Peer) (interface{}, error) {
		conn, err := p.Connect(grpc.WithInsecure(), grpc.WithUnaryInterceptor(message.UnaryClientInterceptor()))
		if err != nil {
			return nil, fmt.Errorf("err: failed to connect to peer(%s) %s\n", p.Addr(), err)
		}
		defer conn.Close()

		s := gw.NewBlockchainServiceClient(conn)
		res, err := s.GetBlockByHeight(ctx, &gw.GetBlockByHeightRequest{BlockHeight: height})
		if err != nil {
			return nil, err
		}

		return res, nil
	}

	res, err := c.invoke(req)
	if err != nil {
		return nil, err
	}

	blk, err := toBlock(res.(*gw.GetBlockResponse).GetBlock())
	if err != nil {
		return nil, err
	}

	return blk, err
}

func (c *Client) GetBlockByRange(ctx context.Context, start uint32, end uint32) ([]*block.Block, error) {
	req := func(p *peer.Peer) (interface{}, error) {
		conn, err := p.Connect(grpc.WithInsecure(), grpc.WithUnaryInterceptor(message.UnaryClientInterceptor()))
		if err != nil {
			return nil, fmt.Errorf("err: failed to connect to peer(%s) %s\n", p.Addr(), err)
		}
		defer conn.Close()

		s := gw.NewBlockchainServiceClient(conn)
		res, err := s.GetBlocks(ctx, &gw.GetBlocksRequest{Offset: start, Limit: end - start + 1})
		if err != nil {
			return nil, err
		}

		return res, err
	}

	res, err := c.invoke(req)
	if err != nil {
		return nil, err
	}

	pbBlks := res.(*gw.GetBlocksResponse).GetBlocks()
	blks := make([]*block.Block, 0, len(pbBlks))
	for _, pbBlk := range pbBlks {
		blk, err := toBlock(pbBlk)
		if err != nil {
			return nil, err
		}
		blks = append(blks, blk)
	}

	return blks, nil
}

func (c *Client) GetTransactionByHash(ctx context.Context, hash [32]byte) (*transaction.Transaction, error) {
	req := func(p *peer.Peer) (interface{}, error) {
		conn, err := p.Connect(grpc.WithInsecure(), grpc.WithUnaryInterceptor(message.UnaryClientInterceptor()))
		if err != nil {
			return nil, fmt.Errorf("err: failed to connect to peer(%s) %s\n", p.Addr(), err)
		}
		defer conn.Close()

		s := gw.NewBlockchainServiceClient(conn)
		res, err := s.GetTransactionByHash(ctx, &gw.GetTransactionByHashRequest{TxHash: fmt.Sprintf("%x", hash)})
		if err != nil {
			return nil, err
		}

		return res, err
	}

	res, err := c.invoke(req)
	if err != nil {
		return nil, err
	}

	pbTx := res.(*gw.GetTransactionResponse).GetTransaction()
	tx, err := toTransaction(pbTx)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (c *Client) GetAccount(ctx context.Context, addr string) (*account.Account, error) {
	req := func(p *peer.Peer) (interface{}, error) {
		conn, err := p.Connect(grpc.WithInsecure(), grpc.WithUnaryInterceptor(message.UnaryClientInterceptor()))
		if err != nil {
			return nil, fmt.Errorf("err: failed to connect to peer(%s) %s\n", p.Addr(), err)
		}
		defer conn.Close()

		s := gw.NewBlockchainServiceClient(conn)
		res, err := s.GetAccount(ctx, &gw.GetAccountRequest{Addr: addr})
		if err != nil {
			return nil, err
		}

		return res, err
	}

	res, err := c.invoke(req)
	if err != nil {
		return nil, err
	}

	pbAcc := res.(*gw.GetAccountResponse).GetAccount()
	if err != nil {
		return nil, err
	}

	return toAccount(pbAcc)
}

func (c *Client) target(addr string) (*peer.Peer, error) {
	for _, p := range c.neighbours {
		if p.Addr() == addr {
			return p, nil
		}
	}

	return nil, errors.New("error: peer not connected")
}

func (c *Client) propagate(reqFunc withoutReturnFn) {
	wg := &sync.WaitGroup{}

	for _, p := range c.neighbours {
		go func(p *peer.Peer) {
			wg.Add(1)
			defer wg.Done()

			if err := reqFunc(p); err != nil {
				p.FailCount += 1
				log.Printf("error: invoke rpc method failed. endpoint=%s. err=%s.\n", p.Addr(), err)
			}
		}(p)
	}

	wg.Wait()
	c.tidyUp()
	c.sortNeghbours()
}

//invoke makes requests to each of the neighbours registered.
//it returns the first response
//it is intended to this function when making requests where the origin of the response is not inportant.
func (c *Client) invoke(reqFunc withReturnFn) (interface{}, error) {
	ch := make(chan interface{})
	errList := make([]error, 0)
	maxErrCnt := len(c.neighbours)

	for _, p := range c.neighbours {
		go func(p *peer.Peer) {
			res, err := reqFunc(p)
			if err != nil {
				p.FailCount += 1
				errList = append(errList, err)
				return
			}
			ch <- res
		}(p)
	}

	for {
		select {
		case r := <-ch:
			return r, nil
		default:
			if len(errList) >= maxErrCnt {
				return nil, errors.New(fmt.Sprintf("error: failed to invoke rpc method. errs=%+v", errList))
			}
		}
	}
}
