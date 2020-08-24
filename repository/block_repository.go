package repository

import (
	"go_chain/block"
)

type blockRepository struct{}

func newBlockRepository() *blockRepository {
	return &blockRepository{}
}

func (br *blockRepository) Get(hash string) (*block.Block, error) {
	return &block.Block{}, nil
}

func (br *blockRepository) GetAll() ([]*block.Block, error) {
	return []*block.Block{}, nil
}
