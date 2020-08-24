package repository

import (
	"go_chain/block"
	"sync"
)

var repositories *Repositories
var once sync.Once

type Repositories struct {
	BlockRepository IBlockRepostitory
}

type IBlockRepostitory interface {
	Get(hash string) (*block.Block, error)
	GetAll() ([]*block.Block, error)
}

func New() *Repositories {
	once.Do(func() { repositories = &Repositories{newBlockRepository()} })

	return repositories
}
