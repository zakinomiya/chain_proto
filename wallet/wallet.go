package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"math/big"
)

type Wallet struct {
	*ecdsa.PrivateKey
	*ecdsa.PublicKey
}

func New() *Wallet {
	return &Wallet{}
}

func (w *Wallet) Restore(priHexStr string) {
	var pri *ecdsa.PrivateKey
	pri.D, _ = new(big.Int).SetString(priHexStr, 16)
	pri.PublicKey.Curve = elliptic.P256()
	pri.PublicKey.X, pri.PublicKey.Y = pri.PublicKey.Curve.ScalarBaseMult(pri.D.Bytes())
	w.PrivateKey = pri
	w.PublicKey = &pri.PublicKey
}

func (w *Wallet) Sign() [32]byte {
	return [32]byte{}
}
