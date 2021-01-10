package blockchain

import "chain_proto/account"

func (bc *Blockchain) GetAccount(addr string) (*account.Account, error) {
	acc, err := bc.repository.Account.Get(addr)
	if err != nil {
		return nil, err
	}

	return acc, nil
}
