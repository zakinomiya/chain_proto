package repository

import (
	"bytes"
	"encoding/binary"
	"errors"
	"go_chain/block"
	"log"

	"github.com/jmoiron/sqlx"
)

type BlockRepository struct {
	*database
}

type blockModel struct {
	height        uint32
	hash          []byte
	prevBlockHash []byte
	merkleRoot    []byte
	timestamp     uint32
	bits          uint32
	nonce         uint32
	extraNonce    uint32
	transactions  []byte
}

func (bm *blockModel) fromBlock(b *block.Block) {
	transactions := []byte{}
	for _, tx := range b.Transactions {
		transactions = append(transactions, tx.ToBytes()...)
	}

	bm.height = b.Height
	bm.hash = b.Hash[:]
	bm.prevBlockHash = b.PrevBlockHash[:]
	bm.merkleRoot = b.MerkleRoot
	bm.timestamp = b.Timestamp
	bm.bits = b.Bits
	bm.nonce = b.Nonce
	bm.extraNonce = b.ExtraNonce
	bm.transactions = transactions

}

func (bm *blockModel) toBlock() (*block.Block, error) {
	var hash [32]byte
	buf := bytes.NewReader(bm.hash)
	err := binary.Read(buf, binary.BigEndian, &hash)
	if err != nil {
		return nil, err
	}

	return &block.Block{
		Height: bm.height,
		Hash:   hash,
		// Transactions: t,
		ExtraNonce: bm.extraNonce,
	}, nil

}

func (br *BlockRepository) GetBlockByHash(hash string) (*block.Block, error) {
	rows, err := br.find("get_block_by_hash.sql", map[string]interface{}{"hash": hash})
	if err != nil {
		return nil, err
	}

	rows.Next()
	block := block.New()
	if err := br.scanBlock(block, rows); err != nil {
		return nil, err
	}

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
		block := block.New()
		if err := br.scanBlock(block, rows); err != nil {
			log.Printf("debug: Failed to scan block. height=%d\n", block.Height)
			return nil, err
		}
		blocks = append(blocks, block)
	}

	return blocks, nil
}

func (br *BlockRepository) GetLatestBlocks(num uint32) ([]*block.Block, error) {
	rows, err := br.find("get_latest_blocks.sql", map[string]interface{}{"num": num})
	if err != nil {
		return nil, err
	}

	blocks := []*block.Block{}
	for rows.Next() {
		block := block.New()
		if err := br.scanBlock(block, rows); err != nil {
			log.Printf("debug: Failed to scan block. height=%d\n", block.Height)
			return nil, err
		}
		blocks = append(blocks, block)
	}

	return blocks, nil
}

func (br *BlockRepository) Insert(b *block.Block) error {
	bm := &blockModel{}
	bm.fromBlock(b)
	return br.exec("insert_block.sql", bm)
}

func (br *BlockRepository) scanBlock(block *block.Block, rows *sqlx.Rows) error {
	if err := rows.StructScan(&block); err != nil {
		log.Println("error:", err)
		return err
	}
	return nil
}
