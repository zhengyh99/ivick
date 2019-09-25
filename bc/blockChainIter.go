package bc

//定义迭代器结构体
type BlockChainIter struct {
	BC          *BlockChain //要迭代的区块链
	CurentBlock []byte      //当前区块hash
}

//迭代器下移
func (iter *BlockChainIter) Next() Block {

	block, err := iter.BC.GetBlock(iter.CurentBlock)
	if err != nil {
		panic(err)
	}
	iter.CurentBlock = block.PrivHash
	return block

}

//后面是否还有数据
func (iter *BlockChainIter) HasNext() bool {
	_, err := iter.BC.GetBlock(iter.CurentBlock)
	if err != nil {
		return false
	}
	return true

}
