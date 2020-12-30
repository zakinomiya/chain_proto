package blockchain

import (
	"fmt"
	"chain_proto/block"
	"chain_proto/db/repository"
	"log"
	"os"
	"sync"
)

// BlockchainEvents represents varieties of events the blockchain invokes
// Subscribers should register a channel which events will be streamed into.
type BlockchainEvents string

const (
	// NewBlock event is invokes when a new block is added to the chain
	NewBlock BlockchainEvents = "NEW_BLOCK"
)

// Blockchain is a struct of the chain
type Blockchain struct {
	mutex         sync.Mutex
	chainID       uint32
	subscriptions map[string]chan BlockchainEvents
	height        uint32
	blocks        []*block.Block
	repository    *repository.Repository
}

var blockchain *Blockchain
var once sync.Once

// New returns a new blockchain
func New(chainID uint32, repository *repository.Repository) *Blockchain {
	blockchain = &Blockchain{
		chainID:       chainID,
		repository:    repository,
		subscriptions: make(map[string]chan BlockchainEvents),
	}
	return blockchain
}

func initializeBlockchain() error {
	b, err := blockchain.repository.Block.GetLatestBlock()
	if err != nil {
		log.Println("error: Failed to initialise blockchain")
		return err
	}

	if b == nil {
		log.Println("info: No blocks found in the db. Creating the genesis block")
		genesis, err := block.NewGenesisBlock()
		if err != nil {
			return err
		}

		if !blockchain.AddBlock(genesis) {
			return fmt.Errorf("error: failed to add the genesis block")
		}
		return nil
	}

	log.Println("info: Block record found in the db. Restoring the blockchain")
	blockchain.height = b.Height
	blockchain.ReplaceBlocks([]*block.Block{b})
	return nil
}

// ServiceName returns its service name
func (bc *Blockchain) ServiceName() string {
	return "Blockchain"
}

// Start starts intialises and starts the blockchain
// There must be only one blockchain during the runtime.
func (bc *Blockchain) Start() error {
	once.Do(func() {
		if err := initializeBlockchain(); err != nil {
			log.Println("error:", err)
			log.Println("Exiting the process...")
			os.Exit(1)
		}
	})
	return nil
}

// Stop stops the blockchain gracefully.
// TODO inplement Stop
func (bc *Blockchain) Stop() {
	return
}

// CurrentBlockHeight returns the height of the latest block in the chain.
// The value does not necessarilly match the length of Blockchain.blocks.
func (bc *Blockchain) CurrentBlockHeight() uint32 {
	return bc.LatestBlock().Height
}

// LatestBlock returns the block at the last index of the Blockchain.blocks.
func (bc *Blockchain) LatestBlock() *block.Block {
	return bc.blocks[len(bc.blocks)-1]
}

// Difficulty returns the next mining target.
// TODO implement Difficulty
func (bc *Blockchain) Difficulty() uint32 {
	return 5
}

// Subscribe regeisters a new channel to the subscription.
// Blockchain events are streaming through the channel
func (bc *Blockchain) Subscribe(key string) <-chan BlockchainEvents {
	ch := make(chan BlockchainEvents)
	bc.subscriptions[key] = ch
	return ch
}

// Unsubscribe delets a subscription with the specified key.
func (bc *Blockchain) Unsubscribe(key string) {
	delete(bc.subscriptions, key)
}

// SendEvent sends events to the subscribed channels.
func (bc *Blockchain) SendEvent(eventName BlockchainEvents) {
	for key, ch := range bc.subscriptions {
		log.Printf("debug: sending event(%s) to the subsctiption(%s)\n", eventName, key)
		ch <- eventName
	}
}
