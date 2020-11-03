package common

import (
	"errors"
	"math/rand"
	"strconv"
)

func IntToByteSlice(b int) []byte {
	return []byte(strconv.Itoa(b))
}

/// Pseudo random uint32
func RandomUint32() uint32 {
	return rand.Uint32()
}

func ReadByteInto32(h []byte) ([32]byte, error) {
	if len(h) != 32 {
		return [32]byte{}, errors.New("byte slice length must be 32")
	}

	var bytes [32]byte
	for i, b := range h {
		bytes[i] = b
	}

	return bytes, nil
}
