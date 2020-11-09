package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"math/big"
)

type Wallet ecdsa.PrivateKey

type Signature struct {
	R *big.Int
	S *big.Int
}

const privKeyByteLength = 32

func New() (*Wallet, error) {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	return (*Wallet)(privKey), nil
}

func (w *Wallet) PrivKeyStr() string {
	return hex.EncodeToString(w.D.Bytes())
}

func (w *Wallet) PubKeyStr() string {
	return hex.EncodeToString(append(w.PublicKey.X.Bytes(), w.PublicKey.Y.Bytes()...))
}

func RestoreWallet(privKeyBytes []byte) (*Wallet, error) {
	if len(privKeyBytes) != privKeyByteLength {
		return nil, errors.New("invalid length of private key")
	}

	priv := &ecdsa.PrivateKey{}
	priv.D = new(big.Int).SetBytes(privKeyBytes)
	priv.PublicKey.Curve = elliptic.P256()
	priv.PublicKey.X, priv.PublicKey.Y = priv.PublicKey.Curve.ScalarBaseMult(privKeyBytes)
	return (*Wallet)(priv), nil
}

func (w *Wallet) Sign(data []byte) (*Signature, error) {
	r, s, err := ecdsa.Sign(rand.Reader, (*ecdsa.PrivateKey)(w), data)
	if err != nil {
		return nil, err
	}

	return &Signature{
		R: r,
		S: s,
	}, nil
}

func (w *Wallet) Verify(signature *Signature, data []byte) bool {
	return ecdsa.Verify(&w.PublicKey, data, signature.R, signature.S)
}
