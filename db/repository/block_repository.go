package repository

import (
	"errors"
	"go_chain/block"
	"go_chain/db/models"
	"log"

	"github.com/jmoiron/sqlx"
)

type BlockRepository struct {
	*commands
}

func (br *BlockRepository) GetBlockByHash(hash string) (*block.Block, error) {
	rows, err := br.find("get_block_by_hash.sql", map[string]interface{}{"hash": hash})
	if err != nil {
		return nil, err
	}

	rows.Next()
	bm := &models.BlockModel{}
	if err := br.scanBlock(bm, rows); err != nil {
		return nil, err
	}

	block := &block.Block{}
	bm.ToBlock(block)
	return block, nil
}

func (br *BlockRepository) GetBlocksByRange(start uint32, end uint32) ([]*block.Block, error) {
	if start > end {
		return nil, errors.New("start height should be less than or equal to end height")
	}

	rows, err := br.find("get_blocks_by_range.sql", map[string]interface{}{"start": start, "end": end})
	if err != nil {
		return nil, err
	}

	blocks := []*block.Block{}
	for rows.Next() {
		bm := &models.BlockModel{}
		if err := br.scanBlock(bm, rows); err != nil {
			log.Printf("debug: Failed to scan block. height=%d\n", bm.Height)
			return nil, err
		}
		block := &block.Block{}
		bm.ToBlock(block)
		blocks = append(blocks, block)
	}
	return blocks, nil
}

func (br *BlockRepository) GetLatestBlock() (*block.Block, error) {
	rows, err := br.find("get_latest_block.sql", nil)
	if err != nil {
		return nil, err
	}

	block := block.New(0, 0, [32]byte{}, nil)
	if rows.Next() {
		bm := &models.BlockModel{}
		if err = br.scanBlock(bm, rows); err != nil {
			log.Printf("debug: Failed to scan block. height=%d\n", bm.Height)
			return nil, err
		}
		bm.ToBlock(block)
		log.Printf("info: latest block height=%d", block.Height)
	} else {
		return nil, nil
	}

	return block, nil
}

func (br BlockRepository) Insert(b *block.Block) error {
	bm := &models.BlockModel{}
	bm.FromBlock(b)
	return r.exec("insert_block.sql", bm)
}

func (br BlockRepository) scanBlock(bm *models.BlockModel, rows *sqlx.Rows) error {
	if err := rows.StructScan(&bm); err != nil {
		log.Println("error:", err)
		return err
	}
	return nil
}
