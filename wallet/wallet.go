package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"math/big"
	"strings"

	"github.com/btcsuite/btcutil/base58"
)

type Wallet struct {
	*ecdsa.PrivateKey
	address string
}

type Signature struct {
	R *big.Int
	S *big.Int
}

func RestoreSignature(sig string) (*Signature, error) {
	rs := strings.Split(sig, "RS")

	if len(rs) != 2 {
		return nil, errors.New("Invalid form of signature string")
	}

	r, ok := new(big.Int).SetString(rs[0], 10)
	if ok != true {
		return nil, errors.New("Invalid form of signature string")
	}

	s, ok := new(big.Int).SetString(rs[1], 10)
	if ok != true {
		return nil, errors.New("Invalid form of siganture string")
	}

	return &Signature{
		R: r,
		S: s,
	}, nil
}

func (s *Signature) String() string {
	return s.R.String() + "RS" + s.S.String()
}

const privKeyByteLength = 32

func New() (*Wallet, error) {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	w := &Wallet{PrivateKey: privKey}
	w.pubKeyToAddr()
	return w, nil
}

func (w *Wallet) pubKeyToAddr() {
	s := sha256.New()
	s.Write(w.PublicKey.X.Bytes())
	s.Write(w.PublicKey.Y.Bytes())

	w.address = base58.Encode(s.Sum([]byte{}))
}

func (w *Wallet) Address() string {
	return w.address
}

func (w *Wallet) PrivKeyStr() string {
	return hex.EncodeToString(w.D.Bytes())
}

func (w *Wallet) PubKeyBytes() []byte {
	return append(w.PublicKey.X.Bytes(), w.PublicKey.Y.Bytes()...)
}

func (w *Wallet) PubKeyStr() string {
	x := hex.EncodeToString(w.PublicKey.X.Bytes())
	y := hex.EncodeToString(w.PublicKey.Y.Bytes())
	return x + y
}

func RestoreWallet(privKeyBytes []byte) (*Wallet, error) {
	if len(privKeyBytes) != privKeyByteLength {
		return nil, errors.New("invalid length of private key")
	}

	priv := &ecdsa.PrivateKey{}
	priv.D = new(big.Int).SetBytes(privKeyBytes)
	priv.PublicKey.Curve = elliptic.P256()
	priv.PublicKey.X, priv.PublicKey.Y = priv.PublicKey.Curve.ScalarBaseMult(privKeyBytes)

	w := &Wallet{
		PrivateKey: priv,
	}
	w.pubKeyToAddr()
	return w, nil
}

func (w *Wallet) Sign(data []byte) (*Signature, error) {
	r, s, err := ecdsa.Sign(rand.Reader, w.PrivateKey, data)
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
