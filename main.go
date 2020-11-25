package main

import (
	"encoding/hex"
	"fmt"
)

func main() {
	h, _ := hex.DecodeString("000007cac0b2b4bfb9117d00a6a26944871e1fa903dbfb100e61171150f43534")
	fmt.Println(h)
}
