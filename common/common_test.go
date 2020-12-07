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

func TestIsValidPort(t *testing.T) {
	ports := []string{
		"0", "", "1000", "9000", "9001", "-9002", "10001", "65535", "65536", "1000000", "hello",
	}
	expected := []bool{false, false, false, true, true, false, true, true, false, false}
	results := []bool{}

	for _, p := range ports {
		results = append(results, IsValidPort(p))
	}

	for i, r := range results {
		if expected[i] != r {
			t.Errorf("FAIL: port=%s should be %t but %t", ports[i], expected[i], r)
		}
	}
}
