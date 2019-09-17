package main

import (
	"fmt"
	_ "fmt"
	"zoin/boltUse"
)

func main() {
	// blockChain := bc.NewBlockChain()
	// blockChain.AddBlock("block1")
	// blockChain.AddBlock("block2")
	// for i, block := range blockChain.Blocks {
	// 	fmt.Printf("=====当前区块高度 %d\n", i)
	// 	fmt.Printf("block prevhash: %x \n", block.PrivHash)
	// 	fmt.Printf("block hash: %x \n", block.Hash)
	// 	fmt.Printf("block data %s \n", block.Data)
	// }
	db := boltUse.OpenBoltDB("testDB.db")
	//添加数据
	if err := db.Put([]byte("NiHao"), []byte("你好")); err != nil {
		fmt.Println("db put error:", err)
	}
	if err := db.Put([]byte("shijie"), []byte("世界")); err != nil {
		fmt.Println("db put error:", err)
	}
	//依键得值
	v := db.GET([]byte("NiHao"))
	fmt.Printf("Key:nihao 's value:%s\n", v)
	//遍历默认桶中的键值对
	datas := db.GetAll("")
	for k, v := range datas {
		if value, ok := v.([][]byte); ok { //数据断言
			fmt.Printf("第%v个数据：, [%s] is %s\n", k, value[0], value[1])
		}

	}
	db.Close()

}
