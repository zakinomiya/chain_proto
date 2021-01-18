package wallet

import (
	"encoding/base64"
	"errors"
	"log"
	"math/big"
)

type Signature struct {
	R *big.Int
	S *big.Int
}

func DecodeSigString(sigString string) (*Signature, error) {
	rbuf, err := base64.StdEncoding.DecodeString(sigString)
	if err != nil {
		log.Println("debug: failed to the decode signature string:", sigString)
		return nil, err
	}

	if len(rbuf) != 64 {
		log.Printf("debug: signature string length is invalid. should be 64 but %d\n", len(rbuf))
		return nil, errors.New("error: invalid signature string")
	}

	r := new(big.Int).SetBytes(rbuf[:32])
	s := new(big.Int).SetBytes(rbuf[32:])

	return &Signature{
		R: r,
		S: s,
	}, nil
}

func (s *Signature) String() string {
	return base64.StdEncoding.EncodeToString(append(s.R.Bytes(), s.S.Bytes()...))
}
