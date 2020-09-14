package repository

import (
	"go_chain/block"
	"sync"
)

var repository *Repository
var once sync.Once

type Repository struct {
	BlockRepository IBlockRepostitory
}

type IBlockRepostitory interface {
	Get(hash string) (*block.Block, error)
	GetAll() ([]*block.Block, error)
}

func New() *Repository {
	once.Do(func() { repository = &Repository{newBlockRepository()} })

	return repository
}
