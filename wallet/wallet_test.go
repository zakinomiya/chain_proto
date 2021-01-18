package wallet

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"
)

func TestWallet(t *testing.T) {
	w, err := New()
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%x", w.D.Bytes())
	hash := sha256.Sum256([]byte("hello world"))
	sig, err := w.Sign(hash[:])
	if err != nil {
		t.Fatal(err)
	}

	if w.Verify(sig, hash[:]) {
		t.Log("Signature verified")
		t.Logf("r=%x, s=%x", sig.R.Bytes(), sig.S.Bytes())
	} else {
		t.Errorf("Invalid signature")
	}
}

const (
	rHex       = "fa81a8cecd19ce29d78beb2ba4b13d6bb1995daa5bccda1ccf71f86ddd15d5f5"
	sHex       = "d75dce035b4d40eb767a124257f03e7970a4a801bda9844b3fabdc2c3142af3f"
	privKeyHex = "771099eb09466a7e2e7f1c8a8087b1b10b970a4e1758f70b7596a52394309863"
)

func TestVerify(t *testing.T) {
	privKeyBytes := []byte(privKeyHex)
	privKey := make([]byte, 32)
	_, err := hex.Decode(privKey, privKeyBytes)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	w, err := RestoreWallet(fmt.Sprintf("%x", privKey))
	if err != nil {
		t.Fatal(err)
	}

	hash := sha256.Sum256([]byte("hello world"))
	sig := &Signature{
		R: new(big.Int).SetBytes(decodeHex(rHex)),
		S: new(big.Int).SetBytes(decodeHex(sHex)),
	}

	if w.Verify(sig, hash[:]) {
		t.Log("Signature verified")
		t.Logf("r=%x, s=%x", sig.R.Bytes(), sig.S.Bytes())
	} else {
		t.Errorf("Invalid signature")
	}
}

func decodeHex(str string) []byte {
	buf := make([]byte, 32)
	hex.Decode(buf, []byte(str))
	return buf
}

func TestDecodeSigString(t *testing.T) {
	sigStr := "b27SEVIjobIjMjwk94Y/ivZjlWXuCo9G8c2td/NSafeazzVdF4LfHWiQqW1xIJQyhu/OXrcgcuSY4fDu8K2ZDQ=="

	sig, err := DecodeSigString(sigStr)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Logf("r=%s, s=%s", sig.R.String(), sig.S.String())
}
