package repository

import (
	"errors"
	"go_chain/block"
	"go_chain/common"
	"log"

	"github.com/jmoiron/sqlx"
)

type BlockModel struct {
	Height        uint32 `db:"height"`
	Hash          []byte `db:"hash"`
	PrevBlockHash []byte `db:"prevBlockHash"`
	MerkleRoot    []byte `db:"merkleRoot"`
	Timestamp     int64  `db:"timestamp"`
	Bits          uint32 `db:"bits"`
	Nonce         uint32 `db:"nonce"`
	ExtraNonce    uint32 `db:"extraNonce"`
	TxCount       int    `db:"txCount"`
	Transactions  string `db:"transactions"`
}

func (r *Repository) fromBlock(b *block.Block, bm *BlockModel) error {
	var transactions string
	for i, tx := range b.Transactions {
		if i != 0 {
			transactions += "!!!"
		}
		transactions += tx.TxHashStr()
	}

	bm.Height = b.Height
	bm.Hash = b.Hash[:]
	bm.PrevBlockHash = b.PrevBlockHash[:]
	bm.MerkleRoot = b.MerkleRoot
	bm.Timestamp = b.Timestamp
	bm.Bits = b.Bits
	bm.Nonce = b.Nonce
	bm.ExtraNonce = b.ExtraNonce
	bm.TxCount = len(b.Transactions)
	bm.Transactions = transactions

	return nil
}

func (r *Repository) toBlock(bm *BlockModel, b *block.Block) error {
	hash := common.ReadByteInto32(bm.Hash)

	transactions, err := r.GetTxByBlockHash(hash)
	if err != nil {
		return err
	}

	b.Hash = hash
	b.Height = bm.Height
	b.MerkleRoot = bm.MerkleRoot
	b.PrevBlockHash = common.ReadByteInto32(bm.PrevBlockHash)
	b.Timestamp = bm.Timestamp
	b.Bits = bm.Bits
	b.ExtraNonce = bm.ExtraNonce
	b.Nonce = bm.Nonce
	b.Transactions = transactions

	return nil
}

func (r *Repository) GetBlockByHash(hash string) (*block.Block, error) {
	rows, err := r.find("get_block_by_hash.sql", map[string]interface{}{"hash": hash})
	if err != nil {
		return nil, err
	}

	rows.Next()
	bm := &BlockModel{}
	if err := r.scanBlock(bm, rows); err != nil {
		return nil, err
	}

	block := &block.Block{}
	r.toBlock(bm, block)
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
		bm := &BlockModel{}
		if err := r.scanBlock(bm, rows); err != nil {
			log.Printf("debug: Failed to scan block. height=%d\n", bm.Height)
			return nil, err
		}
		block := &block.Block{}
		r.toBlock(bm, block)
		blocks = append(blocks, block)
	}
	return blocks, nil
}

func (r *Repository) GetLatestBlock() (*block.Block, error) {
	rows, err := r.find("get_latest_block.sql", nil)
	if err != nil {
		return nil, err
	}

	block := block.New(0, 0, [32]byte{}, nil)
	if rows.Next() {
		bm := &BlockModel{}
		if err := r.scanBlock(bm, rows); err != nil {
			log.Printf("debug: Failed to scan block. height=%d\n", bm.Height)
			return nil, err
		}
		log.Printf("debug: latest block height=%d", block.Height)
		r.toBlock(bm, block)
	} else {
		return nil, nil
	}

	return block, nil
}

func (r *Repository) Insert(b *block.Block) error {
	bm := &BlockModel{}
	r.fromBlock(b, bm)
	return r.exec("insert_block.sql", bm)
}

func (r *Repository) scanBlock(bm *BlockModel, rows *sqlx.Rows) error {
	if err := rows.StructScan(&bm); err != nil {
		log.Println("error:", err)
		return err
	}
	return nil
}
