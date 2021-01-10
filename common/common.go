package common

import (
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

func IntToByteSlice(b int) []byte {
	return []byte(strconv.Itoa(b))
}

/// Pseudo random uint32
func RandomUint32() uint32 {
	return rand.Uint32()
}

//  Read first 32 bytes of the
func ReadByteInto32(h []byte) [32]byte {
	var bytes [32]byte
	for i, b := range h[0:32] {
		bytes[i] = b
	}

	return bytes
}

// Timestamp returns a 64-bit integer
// Yup. I will live for billions of years.
func Timestamp() int64 {
	return time.Now().Unix()
}

// ToDecimal returns decimal type from a string.
// Caller may specify prefix to be removed from the string.
func ToDecimal(decStr string, prefix string) (decimal.Decimal, error) {
	d := strings.TrimPrefix(decStr, prefix)
	return decimal.NewFromString(d)
}
