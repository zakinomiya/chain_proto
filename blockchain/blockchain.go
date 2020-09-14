package blockchain

import (
	"context"
	"fmt"
	"go_chain/block"
	"go_chain/config"
	"go_chain/repository"
	"go_chain/utils"
	"log"
	"os"
	"sync"
)

type BlockchainInterface interface {
	Height() uint32
	Difficulty() uint8
	AddBlock(block *block.Block)
}
type Blockchain struct {
	height       uint32
	blocks       []*block.Block
	repositories *repository.Repository
	difficulty   uint8
}

var blockchain *Blockchain
var once sync.Once

func New(conf *config.ConfigSettings) *Blockchain {
	// TODO initialise db based on config

	once.Do(func() {
		if err := initializeBlockchain(); err != nil {
			log.Println(err)
			os.Exit(1)
		}
	})
	return blockchain
}

func initializeBlockchain() error {
	blockchain = &Blockchain{
		repositories: repository.New(),
	}

	blocks, err := blockchain.repositories.BlockRepository.GetAll()

	if err != nil {
		log.Fatalln("Failed to initialise blockchain")
		return err
	}

	if len(blocks) == 0 {
		log.Println("No blocks found in the db. Creating the genesis block")
		genesis := utils.NewGenesisBlock()
		blockchain.AddBlock(genesis)
		return nil
	}

	log.Println("Block record found in the db. Restoring the blockchain")
	blockchain.height = uint32(len(blocks))
	blockchain.ReplaceBlocks(blocks)
	return nil
}

func (BC *Blockchain) ServiceName() string {
	return "Blockchain"
}

func (bc *Blockchain) Start() error {
	return nil
}

func (bc *Blockchain) Stop() {
	return
}

func (bc *Blockchain) Blocks() []*block.Block {
	return bc.blocks
}

func (bc *Blockchain) Difficulty() uint8 {
	return bc.difficulty
}

func (bc *Blockchain) Height() uint32 {
	return bc.height
}
func (bc *Blockchain) SetHeight(height uint32) {
	bc.height = height
}
func (bc *Blockchain) ReplaceBlocks(blocks []*block.Block) {
	bc.blocks = blocks
}

func (bc *Blockchain) AddBlock(block *block.Block) {
	log.Printf("Adding new block: %#v \n", block)
	bc.blocks = append(bc.blocks, block)
}

func (bc *Blockchain) Mining(block *block.Block) {
	ctx := context.Background()
	// ctxP, cancel := context.WithCancel(ctx)

	for i := 0; i < 1000; i++ {
		// go block.Work(ctxP, 3)
	}

	// cancel()
	fmt.Println("Canceled")
	log.Printf("block hash is %x", ctx.Done())
}
