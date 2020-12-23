package models

import (
	"encoding/hex"
	"fmt"
	"go_chain/block"
	"go_chain/common"
)

type BlockModel struct {
	Height        uint32 `db:"height"`
	Hash          string `db:"hash"`
	PrevBlockHash string `db:"prevBlockHash"`
	MerkleRoot    string `db:"merkleRoot"`
	Timestamp     int64  `db:"timestamp"`
	Bits          uint32 `db:"bits"`
	Nonce         uint32 `db:"nonce"`
	ExtraNonce    uint32 `db:"extraNonce"`
	TxCount       int    `db:"txCount"`
	Transactions  string `db:"transactions"`
}

func (bm *BlockModel) FromBlock(b *block.Block) error {
	var transactions string
	for i, tx := range b.Transactions {
		if i != 0 {
			transactions += "!!!"
		}
		transactions += tx.TxHashStr()
	}

	bm.Height = b.Height
	bm.MerkleRoot = fmt.Sprintf("%x", b.MerkleRoot)
	bm.Timestamp = b.Timestamp
	bm.Bits = b.Bits
	bm.Nonce = b.Nonce
	bm.ExtraNonce = b.ExtraNonce
	bm.TxCount = len(b.Transactions)
	bm.Transactions = transactions
	bm.Hash = fmt.Sprintf("%x", b.Hash)
	bm.PrevBlockHash = fmt.Sprintf("%x", b.PrevBlockHash)

	return nil
}

func (bm *BlockModel) ToBlock(b *block.Block) error {
	h, err := hex.DecodeString(bm.Hash)
	if err != nil {
		return err
	}
	merkleRoot, err := hex.DecodeString(bm.MerkleRoot)
	if err != nil {
		return err
	}

	prevBlockHash, err := hex.DecodeString(bm.PrevBlockHash)

	hash := common.ReadByteInto32(h)
	transactions, err := r.GetTxByBlockHash(hash)
	if err != nil {
		return err
	}

	b.Hash = hash
	b.Height = bm.Height
	b.MerkleRoot = merkleRoot
	b.PrevBlockHash = common.ReadByteInto32(prevBlockHash)
	b.Timestamp = bm.Timestamp
	b.Bits = bm.Bits
	b.ExtraNonce = bm.ExtraNonce
	b.Nonce = bm.Nonce
	b.Transactions = transactions

	return nil
}
