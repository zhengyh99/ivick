package main

import (
	"fmt"
	"zoin/bc"
)

func main() {
	blockChain := bc.NewBlockChain()
	blockChain.AddBlock("block1")
	blockChain.AddBlock("block2")
	for i, block := range blockChain.Blocks {
		fmt.Printf("=====当前区块高度 %d\n", i)
		fmt.Printf("block prevhash: %x \n", block.PrivHash)
		fmt.Printf("block hash: %x \n", block.Hash)
		fmt.Printf("block data %s \n", block.Data)
	}

}
