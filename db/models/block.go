package models

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"go_chain/block"
	"go_chain/common"
	"go_chain/transaction"
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
	transactions, err := json.Marshal(b.Transactions)
	if err != nil {
		return err
	}

	bm.Height = b.Height
	bm.MerkleRoot = fmt.Sprintf("%x", b.MerkleRoot)
	bm.Timestamp = b.Timestamp
	bm.Bits = b.Bits
	bm.Nonce = b.Nonce
	bm.ExtraNonce = b.ExtraNonce
	bm.TxCount = len(b.Transactions)
	bm.Transactions = string(transactions)
	bm.Hash = fmt.Sprintf("%x", b.Hash)
	bm.PrevBlockHash = fmt.Sprintf("%x", b.PrevBlockHash)

	return nil
}

func (bm *BlockModel) ToBlock() (*block.Block, error) {
	h, err := hex.DecodeString(bm.Hash)
	if err != nil {
		return nil, err
	}
	merkleRoot, err := hex.DecodeString(bm.MerkleRoot)
	if err != nil {
		return nil, err
	}

	prevBlockHash, err := hex.DecodeString(bm.PrevBlockHash)

	hash := common.ReadByteInto32(h)

	var transactions []*transaction.Transaction
	if err := json.Unmarshal([]byte(bm.Transactions), &transactions); err != nil {
		return nil, err
	}

	b := block.New(bm.Height, bm.Bits, common.ReadByteInto32(prevBlockHash), transactions)
	b.Hash = hash
	b.MerkleRoot = merkleRoot
	b.Timestamp = bm.Timestamp
	b.ExtraNonce = bm.ExtraNonce
	b.Nonce = bm.Nonce

	return b, nil
}
