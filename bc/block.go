package bc

import "crypto/sha256"

type Block struct {
	PrivHash []byte
	Hash     []byte
	Data     []byte
}

func NewBlock(data string, prevBlockHash []byte) (block *Block) {
	block = &Block{
		PrivHash: prevBlockHash,
		Hash:     []byte{},
		Data:     []byte(data),
	}
	block.setHash()
	return
}

func (block *Block) setHash() {
	blockInfo := append(block.PrivHash, block.Data...)
	hash := sha256.Sum256(blockInfo)
	block.Hash = hash[:]
}
func GenesisBlock() *Block {
	return NewBlock("创世块", []byte{})
}
