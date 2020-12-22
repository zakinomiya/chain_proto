package account

import (
	"errors"
	"log"
)

// Account is a model of user in the blockchain.
type Account struct {
	addr    string
	balance uint32
}

// New initialises a new Account struct
// TODO implement Account New. Need to fetch account from the db.
func New(addr string) *Account {
	return &Account{addr: addr, balance: 0}
}

// Addr is a getter for addr
func (a *Account) Addr() string {
	return a.addr
}

// Balance is a getter for balance
func (a *Account) Balance() uint32 {
	return a.balance
}

// Send sends a specified amount of coins from an account to another
func (a *Account) Send(amount uint32, recipient *Account) error {
	if a.balance < amount {
		return errors.New("error: balance not enough")
	}
	a.balance -= amount
	recipient.Receive(amount)
	return nil
}

// Receive sums a specified amount of coins to the current balance
func (a *Account) Receive(amount uint32) {
	log.Printf("info: account(%s) received %d coins\n", a.addr, amount)
	a.balance += amount
}
