package repository

import (
	"chain_proto/account"
	"chain_proto/block"
	"chain_proto/db/models"
	"chain_proto/transaction"
	"database/sql"
	"errors"
	"log"

	"github.com/jmoiron/sqlx"
)

type BlockRepository struct {
	*database
	account     *AccountRepository
	transcation *TxRepository
}

func (br *BlockRepository) GetBlockByHash(hash string) (*block.Block, error) {
	row, err := br.queryRow("get_block_by_hash.sql", map[string]interface{}{"hash": hash})
	if err != nil {
		return nil, err
	}

	bm := &models.BlockModel{}
	if err := row.StructScan(&bm); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return bm.ToBlock()
}

func (br *BlockRepository) GetBlocksByRange(start uint32, end uint32) ([]*block.Block, error) {
	if start > end {
		return nil, errors.New("start height should be less than or equal to end height")
	}

	rows, err := br.queryRows("get_blocks_by_range.sql", map[string]interface{}{"start": start, "end": end})
	if err != nil {
		return nil, err
	}

	blocks := []*block.Block{}
	for rows.Next() {
		bm := &models.BlockModel{}
		if err := rows.StructScan(&bm); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			log.Printf("debug: Failed to scan block. height=%d\n", bm.Height)
			return nil, err
		}
		block, err := bm.ToBlock()
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, block)
	}

	return blocks, nil
}

func (br *BlockRepository) GetLatestBlock() (*block.Block, error) {
	log.Println("debug: action=GetLatestBlock")
	row, err := br.queryRow("get_latest_block.sql", nil)
	if err != nil {
		return nil, err
	}

	bm := &models.BlockModel{}
	if err = row.StructScan(bm); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Printf("debug: Failed to scan block. height=%d\n", bm.Height)
		return nil, err
	}

	log.Printf("info: latest block height=%d", bm.Height)

	return bm.ToBlock()
}

func (br BlockRepository) Insert(b *block.Block) error {
	log.Println("debug: action=BlockRepository.Insert")
	tx, err := br.db.Beginx()
	if err != nil {
		return err
	}

	if err := br.insert(tx, b); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (br *BlockRepository) insert(tx *sqlx.Tx, b *block.Block) error {
	filename := "insert_block.sql"
	processedAccounts, err := br.processTxs(b.Transactions)
	if err != nil {
		return err
	}

	bm := &models.BlockModel{}
	bm.FromBlock(b)

	if err := br.txCommand(tx, filename, bm); err != nil {
		return err
	}
	if err := br.account.bulkInsert(tx, processedAccounts); err != nil {
		return err
	}
	if err := br.transcation.bulkInsert(tx, b.Transactions); err != nil {
		return err
	}

	return nil
}

// processTxs processes the transactions in a block and simulates the blockchain state after all the transactions become valid.
// process failes when sender does not have enough balance.
func (br *BlockRepository) processTxs(txs []*transaction.Transaction) ([]*account.Account, error) {
	log.Println("debug: action=processTxs")
	accountMap := map[string]*account.Account{}

	for _, tx := range txs {
		log.Printf("debug: SenderAddr=%s", tx.SenderAddr)
		sender := accountMap[tx.SenderAddr]
		if sender == nil && tx.TxType != "coinbase" {
			s, err := br.account.GetAccount(tx.SenderAddr)
			if err != nil {
				return nil, err
			}
			sender = s
			accountMap[tx.SenderAddr] = sender
		}

		for _, output := range tx.Outs {
			log.Printf("debug: RecipientAddr=%s", output.RecipientAddr)
			recipient := accountMap[output.RecipientAddr]
			if recipient == nil {
				r, err := br.account.GetAccount(output.RecipientAddr)
				if err != nil {
					return nil, err
				}
				recipient = r
				accountMap[output.RecipientAddr] = recipient
			}
			if tx.TxType == "coinbase" {
				log.Printf("debug: sending coinbase tx to %s\n", recipient.Addr)
				recipient.Receive(output.Value)
				continue
			}

			if err := sender.Send(output.Value, recipient); err != nil {
				log.Printf("error: failed to send amount from %s to %s\n", sender.Addr, recipient.Addr)
				return nil, err
			}
		}
	}
	accounts := []*account.Account{}
	for _, account := range accountMap {
		accounts = append(accounts, account)
	}

	log.Printf("debug: action=processTxs. accounts=%+v\n", accounts)
	return accounts, nil
}
