package gateway

import (
	"chain_proto/account"
	"chain_proto/block"
	gw "chain_proto/gateway/gw"
	"chain_proto/peer"
	"chain_proto/transaction"
)

type Blockchain interface {
	GetBlockByHash(hash [32]byte) (*block.Block, error)
	GetBlockByHeight(height uint32) (*block.Block, error)
	GetBlocks(offset uint32, limit uint32) ([]*block.Block, error)
	AddBlock(block *block.Block) error
	LatestBlock() *block.Block
	GetTxsByBlockHash(blockHash [32]byte) ([]*transaction.Transaction, error)
	GetTransactionByHash(hash [32]byte) (*transaction.Transaction, error)
	AddTransaction(tx *transaction.Transaction) error
	GetAccount(addr string) (*account.Account, error)
	AddPeer(addr string, network string) error
	GetPeers() ([]*peer.Peer, error)
}

type BlockchainService struct {
	bc Blockchain
	gw.UnimplementedBlockchainServiceServer
}

func NewBlockchainService(bc Blockchain) *BlockchainService {
	return &BlockchainService{bc: bc}
}
