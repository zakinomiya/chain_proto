package account

import (
	"errors"
	"log"
)

// Account is a model of user in the blockchain.
type Account struct {
	Addr    string
	Balance uint32
}

// New initialises a new Account struct
// TODO implement Account New. Need to fetch account from the db.
func New(addr string) *Account {
	return &Account{Addr: addr, Balance: 0}
}

// Send sends a specified amount of coins from an account to another
func (a *Account) Send(amount uint32, recipient *Account) error {
	if a.Balance < amount {
		return errors.New("error: balance not enough")
	}
	a.Balance -= amount
	recipient.Receive(amount)
	return nil
}

// Receive sums a specified amount of coins to the current balance
func (a *Account) Receive(amount uint32) {
	log.Printf("info: account(%s) received %d coins\n", a.Addr, amount)
	a.Balance += amount
}
