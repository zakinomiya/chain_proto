package common

import (
	"crypto/sha256"
	"fmt"
	"testing"
)

func TestReadByteInto32(t *testing.T) {
	sha := sha256.New()
	sha.Write([]byte("hello"))
	sum := sha.Sum(nil)
	sumStr := fmt.Sprintf("%x", sum)

	s := ReadByteInto32(sum)

	sss := fmt.Sprintf("%x", s[:])

	if sumStr != sss {
		t.Logf("sumStr=%s, sss=%s\n", sumStr, sss)
		t.Fatal("test failed")
	}
	t.Logf("sumStr=%s, sss=%s\n", sumStr, sss)
}
