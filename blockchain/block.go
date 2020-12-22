package blockchain

import (
	"go_chain/account"
	"go_chain/block"
	"go_chain/transaction"
	"log"
)

// Blocks returns Blockchain.blocks
func (bc *Blockchain) Blocks() []*block.Block {
	return bc.blocks
}

// ReplaceBlocks replaces the entire chain with the new one.
func (bc *Blockchain) ReplaceBlocks(blocks []*block.Block) {
	bc.blocks = blocks
}

// AddBlock adds a new block to the chain.
func (bc *Blockchain) AddBlock(block *block.Block) bool {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()

	if !block.Verify() {
		log.Println("info: Refused adding the block")
		return false
	}

	log.Printf("info: Adding new block: %x\n", block.Hash)
	log.Printf("debug: block=%+v\n", block)
	log.Printf("info: Now the length of the chain is %d:\n", bc.LatestBlock().Height)

	if err := bc.repository.Insert(block); err != nil {
		log.Printf("error: failed to insert block to db. err=%v", err)
		return false
	}

	bc.blocks = append(bc.blocks, block)
	bc.SendEvent(NewBlock)
	return true
}

func (bc *Blockchain) processTxs(txs []*transaction.Transaction) (map[string]*account.Account, error) {
	log.Println("debug: processTxs")
	accounts := map[string]*account.Account{}

	for _, tx := range txs {
		sender := accounts[tx.SenderAddr]
		if sender == nil {
			sender = account.New(tx.SenderAddr)
			accounts[tx.SenderAddr] = sender
		}

		for _, output := range tx.Outs {
			recipient := accounts[output.RecipientAddr]
			if recipient == nil {
				recipient = account.New(output.RecipientAddr)
				accounts[output.RecipientAddr] = recipient
			}

			if err := sender.Send(output.Value, recipient); err != nil {
				log.Printf("error: failed to send amount from %s to %s\n", sender.Addr(), recipient.Addr())
				return nil, err
			}
		}
	}

	return accounts, nil
}
