package blockchain

import "chain_proto/account"

// GetAccount returns an account with the given address
func (bc *Blockchain) GetAccount(addr string) (*account.Account, error) {
	acc, err := bc.r.Account.Get(addr)
	if err != nil {
		return nil, err
	}

	return acc, nil
}
