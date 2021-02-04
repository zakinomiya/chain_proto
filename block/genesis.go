package block

import (
	"chain_proto/common"
	"chain_proto/transaction"
	"chain_proto/wallet"
	"encoding/hex"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/go-yaml/yaml"
	"github.com/shopspring/decimal"
)

type genesis struct {
	PrevBlockHash string `yaml:"prevBlockHash"`
	MerkleRoot    string `yaml:"merkleRoot"`
	Timestamp     int64  `yaml:"timestamp"`
	Bits          uint32 `yaml:"bits"`
	Nonce         uint32 `yaml:"nonce"`
	Hash          string `yaml:"hash"`
	ExtraNonce    uint32 `yaml:"extraNonce"`
	Transactions  []struct {
		TxHash     string             `yaml:"txHash"`
		TxType     transaction.TxType `yaml:"txType"`
		TotalValue string             `yaml:"totalValue"`
		Fee        string             `yaml:"fee"`
		SenderAddr string             `yaml:"senderAddr"`
		Timestamp  int64              `yaml:"timestamp"`
		Signature  string             `yaml:"signature"`
		Outs       []struct {
			RecipientAddr string `yaml:"recipientAddr"`
			Value         string `yaml:"value"`
		} `yaml:"outs"`
	} `yaml:"transactions"`
}

func readFromYaml(path string) (*genesis, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	gen := &genesis{}

	if err := yaml.Unmarshal(data, gen); err != nil {
		return nil, err
	}

	log.Printf("debug: gen=%+v", gen)
	return gen, nil
}

func NewGenesisBlock() (*Block, error) {
	gen, err := readFromYaml(filepath.Join(os.Getenv("GOPATH"), "src/chain_proto/config/genesis.yaml"))

	var transactions []*transaction.Transaction
	for _, t := range gen.Transactions {
		txHash, err := hex.DecodeString(t.TxHash)
		if err != nil {
			log.Println("error:", t.TxHash)
			return nil, err
		}

		totalValue, _ := decimal.NewFromString(t.TotalValue)
		fee, _ := decimal.NewFromString(t.Fee)

		tx := transaction.New()
		tx.TxHash = common.ReadByteInto32(txHash)
		tx.TxType = t.TxType
		tx.SenderAddr = t.SenderAddr
		tx.Timestamp = t.Timestamp
		tx.TotalValue = totalValue
		tx.Fee = fee
		sig, err := wallet.DecodeSigString(t.Signature)
		tx.Signature = sig

		var outs []*transaction.Output
		for _, o := range t.Outs {
			value, _ := decimal.NewFromString(o.Value)
			out := &transaction.Output{}
			out.RecipientAddr = o.RecipientAddr
			out.Value = value
			if err != nil {
				return nil, err
			}
			outs = append(outs, out)
		}
		tx.Outs = outs

		transactions = append(transactions, tx)
	}

	hash, err := hex.DecodeString(gen.Hash)
	if err != nil {
		log.Println("error:", gen.Hash)
		return nil, err
	}

	prevBlockHash, err := hex.DecodeString(gen.PrevBlockHash)
	if err != nil {
		log.Println("error:", gen.PrevBlockHash)
		return nil, err
	}

	merklerRoot, err := hex.DecodeString(gen.MerkleRoot)
	if err != nil {
		log.Println("error:", gen.MerkleRoot)
		return nil, err
	}

	b := &Block{
		Height:       0,
		Hash:         common.ReadByteInto32(hash),
		ExtraNonce:   gen.ExtraNonce,
		Transactions: transactions,
		BlockHeader: &BlockHeader{
			PrevBlockHash: common.ReadByteInto32(prevBlockHash),
			MerkleRoot:    merklerRoot,
			Bits:          gen.Bits,
			Nonce:         gen.Nonce,
			Timestamp:     gen.Timestamp,
		},
	}

	b.Hash = b.HashBlock()

	log.Printf("debug: genesis block=%+v\n", b)
	return b, nil
}
