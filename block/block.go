package block

type Block struct {
	hash   string
	amount int
}

func New() *Block {
	return &Block{"some hash", 100}
}

func (block *Block) Hash() string {
	return block.hash
}

func (block *Block) Amount() int {
	return block.amount
}

func (block *Block) SetHash(hash string) *Block {
	block.hash = hash
	return block
}

func (block *Block) SetAmount(amount int) *Block {
	block.amount = amount
	return block
}
