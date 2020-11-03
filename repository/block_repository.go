package repository

import (
	"errors"
	"go_chain/block"
	"go_chain/common"
	"log"

	"github.com/jmoiron/sqlx"
)

type blockModel struct {
	height        uint32
	hash          []byte
	prevBlockHash []byte
	merkleRoot    []byte
	timestamp     uint32
	bits          uint32
	nonce         uint32
	extraNonce    uint32
	txCount       int
	transactions  []string
}

func (r *Repository) fromBlock(b *block.Block, bm *blockModel) error {
	var transactions []string
	for _, tx := range b.Transactions {
		transactions = append(transactions, tx.TxHashStr())
	}

	bm.height = b.Height
	bm.hash = b.Hash[:]
	bm.prevBlockHash = b.PrevBlockHash[:]
	bm.merkleRoot = b.MerkleRoot
	bm.timestamp = b.Timestamp
	bm.bits = b.Bits
	bm.nonce = b.Nonce
	bm.extraNonce = b.ExtraNonce
	bm.txCount = len(b.Transactions)
	bm.transactions = transactions

	return nil
}

func (r *Repository) toBlock(bm *blockModel, b *block.Block) error {
	hash, err := common.ReadByteInto32(bm.hash)
	if err != nil {
		return err
	}

	prevBlockHash, err := common.ReadByteInto32(bm.prevBlockHash)
	if err != nil {
		return err
	}

	transactions, err := r.GetTxByBlockHash(hash)
	if err != nil {
		return err
	}

	b.Hash = hash
	b.Height = bm.height
	b.MerkleRoot = bm.merkleRoot
	b.PrevBlockHash = prevBlockHash
	b.Timestamp = bm.timestamp
	b.Bits = bm.bits
	b.ExtraNonce = bm.extraNonce
	b.Nonce = bm.nonce
	b.Transactions = transactions

	return nil
}

func (r *Repository) GetBlockByHash(hash string) (*block.Block, error) {
	rows, err := r.find("get_block_by_hash.sql", map[string]interface{}{"hash": hash})
	if err != nil {
		return nil, err
	}

	rows.Next()
	block := block.New()
	if err := r.scanBlock(block, rows); err != nil {
		return nil, err
	}

	return block, nil
}

func (r *Repository) GetBlocksByRange(start uint32, end uint32) ([]*block.Block, error) {
	if start > end {
		return nil, errors.New("start height should be less than or equal to end height")
	}

	rows, err := r.find("get_blocks_by_range.sql", map[string]interface{}{"start": start, "end": end})
	if err != nil {
		return nil, err
	}

	blocks := []*block.Block{}
	for rows.Next() {
		block := block.New()
		if err := r.scanBlock(block, rows); err != nil {
			log.Printf("debug: Failed to scan block. height=%d\n", block.Height)
			return nil, err
		}
		blocks = append(blocks, block)
	}

	return blocks, nil
}

func (br *Repository) GetLatestBlocks(num uint32) ([]*block.Block, error) {
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

func (br *Repository) Insert(b *block.Block) error {
	bm := &blockModel{}
	br.fromBlock(b, bm)
	return br.exec("insert_block.sql", bm)
}

func (br *Repository) scanBlock(block *block.Block, rows *sqlx.Rows) error {
	if err := rows.StructScan(&block); err != nil {
		log.Println("error:", err)
		return err
	}
	return nil
}
