package bc

type BlockChain struct {
	Blocks []*Block
}

func NewBlockChain() *BlockChain {
	genesisBlock := GenesisBlock()
	return &BlockChain{
		Blocks: []*Block{genesisBlock},
	}
}

func (bc *BlockChain) AddBlock(data string) {
	lastBlock := bc.Blocks[len(bc.Blocks)-1]
	block := NewBlock(data, lastBlock.Hash)
	bc.Blocks = append(bc.Blocks, block)
}
