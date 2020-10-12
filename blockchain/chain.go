package blockchain

import (
	"go_chain/block"
	"go_chain/config"
	"go_chain/repository"
	"go_chain/transaction"
	"log"
	"os"
	"sync"
)

type BlockchainInterface interface {
	Height() uint32
	GenerateBlock(txs []*transaction.Transaction) *block.Block
	AddBlock(block *block.Block) bool
}

type Blockchain struct {
	lock         sync.RWMutex
	chainID      uint32
	height       uint32
	blocks       []*block.Block
	repositories *repository.Repository
}

var blockchain *Blockchain
var once sync.Once

func New(conf *config.ConfigSettings) *Blockchain {

	once.Do(func() {
		if err := initializeBlockchain(conf.ChainID); err != nil {
			log.Println(err, "Exiting the process...")
			os.Exit(1)
		}
	})
	return blockchain
}

func initializeBlockchain(chainID uint32) error {
	blockchain = &Blockchain{
		chainID:      chainID,
		repositories: repository.New(),
	}

	blocks, err := blockchain.repositories.BlockRepository.GetAll()

	if err != nil {
		log.Fatalln("error: Failed to initialise blockchain")
		return err
	}

	if len(blocks) == 0 {
		log.Println("info: No blocks found in the db. Creating the genesis block")
		// genesis := utils.NewGenesisBlock()
		// // if !blockchain.AddBlock(genesis) {
		// // 	return fmt.Errorf("error: failed to add the genesis block")
		// // }
		return nil
	}

	log.Println("info: Block record found in the db. Restoring the blockchain")
	blockchain.height = uint32(len(blocks))
	blockchain.ReplaceBlocks(blocks)
	return nil
}

func (bc *Blockchain) ServiceName() string {
	return "Blockchain"
}

func (bc *Blockchain) Start() error {
	return nil
}

func (bc *Blockchain) Stop() {
	return
}

func (bc *Blockchain) Height() uint32 {
	return bc.height
}
func (bc *Blockchain) SetHeight(height uint32) {
	bc.height = height
}
