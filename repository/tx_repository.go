package repository

type TxRepository struct {
	*database
}

type txModel struct {
	txHash        []byte
	totalValue    uint32
	fee           uint32
	senderAddr    []byte
	timestamp     uint64
	recipientAddr []byte
	value         uint32
	signature     []byte
}
