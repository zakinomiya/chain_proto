package common

import (
	"math/rand"
	"strconv"
	"time"
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
