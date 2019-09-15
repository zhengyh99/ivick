package bc

//定义区块链
type BlockChain struct {
	Blocks []*Block
}

//创建区块链
func NewBlockChain() *BlockChain {
	genesisBlock := GenesisBlock()
	return &BlockChain{
		Blocks: []*Block{genesisBlock},
	}
}

//向区块链中加入区块
func (bc *BlockChain) AddBlock(data string) {
	lastBlock := bc.Blocks[len(bc.Blocks)-1]
	block := NewBlock(data, lastBlock.Hash)
	bc.Blocks = append(bc.Blocks, block)
}
