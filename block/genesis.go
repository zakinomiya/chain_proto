package block

import (
	"go_chain/transaction"
	"io/ioutil"

	"github.com/go-yaml/yaml"
)

type Genesis struct {
	PrevBlockHash [32]byte `yaml:"prevBlockHash"`
	MerkleRoot    []byte   `yaml:"merkleRoot"`
	Timestamp     uint32   `yaml:"timestamp"`
	Bits          uint32   `yaml:"bits"`
	Nonce         uint32   `yaml:"nonce"`
	Height        uint32   `yaml:"height"`
	Hash          [32]byte `yaml:"hash"`
	ExtraNonce    uint32   `yaml:"extraNonce"`
	Transactions  []struct {
		TxHash     [32]byte `yaml:"txHash"`
		TotalValue uint32   `yaml:"totalValue"`
		Fee        uint32   `yaml:"fee"`
		SenderAddr [32]byte `yaml:"senderAddr"`
		Timestamp  uint64   `yaml:"timestamp"`
		Outs       []struct {
			Index         uint32   `yaml:"index"`
			RecipientAddr [32]byte `yaml:"recipientAddr"`
			Value         uint32   `yaml:"value"`
			Signature     string   `yaml:"signature"`
		} `yaml:"outs"`
	} `yaml:"transactions"`
}

func NewGenesisBlock() (*Block, error) {
	data, err := ioutil.ReadFile("./block/genesis.yaml")
	if err != nil {
		return nil, err
	}

	gen := &Genesis{}

	if err := yaml.Unmarshal(data, gen); err != nil {
		return nil, err
	}

	var transactions []*transaction.Transaction
	for _, t := range gen.Transactions {
		tx := transaction.New()
		tx.TxHash = t.TxHash
		tx.SenderAddr = t.SenderAddr
		tx.Timestamp = t.Timestamp
		tx.TotalValue = t.TotalValue
		tx.Fee = t.Fee

		var outs []*transaction.Output
		for _, o := range t.Outs {
			out := &transaction.Output{}
			out.Index = o.Index
			out.Value = o.Value
			out.Signature = []byte{}
			out.RecipientAddr = o.RecipientAddr
			outs = append(outs, out)
		}

		tx.Outs = outs
		transactions = append(transactions, tx)
	}

	b := &Block{
		Height:       gen.Height,
		Hash:         gen.Hash,
		ExtraNonce:   gen.ExtraNonce,
		Transactions: transactions,
		BlockHeader: &BlockHeader{
			PrevBlockHash: gen.PrevBlockHash,
			MerkleRoot:    gen.MerkleRoot,
			Bits:          gen.Bits,
			Nonce:         gen.Nonce,
			Timestamp:     gen.Timestamp,
		},
	}

	return b, nil
}
