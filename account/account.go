package account

type Account struct {
	addr    [32]byte
	balance uint32
}

func New(addr [32]byte, balance uint32) Account {
	return Account{addr: addr, balance: balance}
}

func (a *Account) Addr() [32]byte {
	return a.addr
}

func (a *Account) Balance() uint32 {
	return a.balance
}
