package wallet

type Wallet struct {
	secKey [32]byte
	pubKey [32]byte
}

func New(secKey, pubKey [32]byte) Wallet {
	return Wallet{secKey: secKey, pubKey: pubKey}
}

func (w *Wallet) Sign() [32]byte {
	return [32]byte{}
}
