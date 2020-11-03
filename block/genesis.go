package block

func NewGenesisBlock() *Block {
	b := New()
	h := NewHeader()
	h.Bits = 5
	h.Nonce = 129
	h.MerkleRoot = []byte("merkle")
	h.PrevBlockHash = [32]byte{}
	// coinbase := NewCoinbase([]byte("This is Minimum Viable Blockchain"), 25)
	b.BlockHeader = h
	b.Height = 1000
	// b.Transactions = []*transaction.Transaction{coinbase}
	b.SetExtranNonce()
	return b
}