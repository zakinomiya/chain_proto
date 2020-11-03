package transaction

func NewCoinbase(pubKey []byte, value uint64) *Transaction {
	tx := New()
	tx.CalcHash()
	return tx
}
