package repository

import (
	"chain_proto/block"
	"chain_proto/transaction"
	"chain_proto/wallet"
	"crypto/sha256"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var r *Repository

const (
	dbPath   = "blockchain_test.db"
	dbDriver = "sqlite3"
	sqlPath  = "../sql"
)

func setup() func() {
	fmt.Println("repository test started")

	if _, err := os.Stat(dbPath); !os.IsNotExist(err) {
		fmt.Println("test db file exists. removing file")
		os.Remove(dbPath)
	}

	repo, err := New(dbPath, dbDriver, sqlPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	r = repo

	return func() {
		fmt.Println("repository test done")
	}
}

func TestMain(m *testing.M) {
	shutdown := setup()
	defer shutdown()

	m.Run()
}

func TestBlockRepository(t *testing.T) {
	t.Log("test started. BlockRepository")
	w1, _ := wallet.New()
	blk1 := block.New(1, 5, [32]byte{}, nil)
	blk1.ExtraNonce = 100000000
	blk1.Hash = sha256.Sum256([]byte("block hash1"))
	blk1.MerkleRoot = []byte("merkle root1")
	blk1.Transactions = []*transaction.Transaction{transaction.NewCoinbase(w1, "25.000")}

	w2, _ := wallet.New()
	blk2 := block.New(2, 5, [32]byte{}, nil)
	blk2.ExtraNonce = 100000000
	blk2.Hash = sha256.Sum256([]byte("block hash2"))
	blk2.MerkleRoot = []byte("merkle root2")
	blk2.Transactions = []*transaction.Transaction{transaction.NewCoinbase(w2, "25.000")}

	t.Run("insert", func(t *testing.T) {
		err := r.Block.Insert(blk1)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log("insert done")
	})

	t.Run("GetByHash", func(t *testing.T) {
		b, err := r.Block.GetByHash(blk1.Hash)
		if err != nil {
			t.Errorf("GetByHash failed. err %+v\n", err)
			return
		}

		assert.Equal(t, blk1, b)
		t.Log("GetByHash done")
	})

	t.Run("GetByHeight", func(t *testing.T) {
		b, err := r.Block.GetByHeight(1)
		if err != nil {
			t.Errorf("GetByHeight failed. err %+v\n", err)
			return
		}

		assert.Equal(t, blk1, b)
		t.Log("GetByHeight done")
	})

	t.Run("GetLatest", func(t *testing.T) {
		b2, err := r.Block.GetLatest()
		if err != nil {
			t.Errorf("GetLatest failed. err %+v\n", err)
			return
		}

		assert.Equal(t, blk2, b2)
		t.Log("GetLatest done")
	})

	t.Run("GetByRange", func(t *testing.T) {
		blks, err := r.Block.GetByRange(0, 2)
		if err != nil {
			t.Errorf("GetByRange failed. err %+v\n", err)
			return
		}
		if len(blks) != 2 {
			t.Errorf("GetByRange falied. err=expected blocks length is 2, but given %d\n", len(blks))
			return
		}

		assert.Equal(t, blk1, blks[0])
		assert.Equal(t, blk2, blks[1])
		t.Log("GetByRange done")
	})
}
