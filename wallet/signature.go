package wallet

import (
	"errors"
	"log"
	"math/big"
	"strings"
)

type Signature struct {
	R *big.Int
	S *big.Int
}

func RestoreSignature(sig string) (*Signature, error) {
	rs := strings.Split(sig, "RS")

	if len(rs) != 2 {
		log.Printf("debug: sig string=%s. rs=%+v\n", sig, rs)
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
