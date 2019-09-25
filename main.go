package main

import (
	"zoin/bc"
)

func main() {

	cli := bc.NewCLI()
	defer cli.Close()
	cli.Run()

	// //测试 序列化
	// block := bc.NewBlock("hahahaha", []byte("helloword"))
	// enBlock := block.Serialize()
	// b, err := bc.DeSerliaze(enBlock)
	// if err != nil {
	// 	fmt.Println("bc deseriaze error:", err)
	// }
	// fmt.Printf("b data:%s,privHash:%s", b.Data, b.PrivHash)

	// //测试blockChain中的 bolt数据操作
	// blockChain := bc.GetBlockChain()
	// blockChain.AddBlock("block3")
	// blockChain.AddBlock("block4")

	// datas := blockChain.GetAll()
	// for k, v := range datas {
	// 	if value, ok := v.([][]byte); ok { //数据断言

	// 		fmt.Printf("第%v个数据：, [%x] is %x\n", k, value[0], value[1])
	// 	}
	// }

	// blockChain.Close()

	//测试 BlockChainIter 迭代器
	// bc := bc.GetBlockChain()
	// for iter := blockChain.Iter(); iter.HasNext(); {
	// 	b := iter.Next()
	// 	fmt.Printf("b.PrivHash：%x\n", b.PrivHash)
	// 	fmt.Printf("b.Hash:%x \n", b.Hash)
	// 	fmt.Printf("b.daga:%s\n", b.Data)
	// 	fmt.Println("==========")
	// }
	// blockChain.Close()

	//测试 boltUse db.go
	// db := boltUse.OpenBoltDB("blockChain.db", "blockChain")
	// ds := db.GET([]byte("tail"))
	// fmt.Printf("ds:%v", ds)
	// // for k, v := range ds {
	// // 	fmt.Println("0000000000000000")
	// // 	if value, ok := v.([][]byte); ok { //数据断言
	// // 		fmt.Printf("第%v个数据：, [%s] is %s\n", k, value[0], value[1])
	// // 	}
	// // }

	// db.Close()

	// db := boltUse.OpenBoltDB("testDB.db")
	// //添加数据
	// if err := db.Put([]byte("NiHao2"), []byte("你好")); err != nil {
	// 	fmt.Println("db put error:", err)
	// }
	// if err := db.Put([]byte("shijie3"), []byte("世界")); err != nil {
	// 	fmt.Println("db put error:", err)
	// }
	// //依键得值
	// v := db.GET([]byte("NiHao"))
	// fmt.Printf("Key:nihao 's value:%s\n", v)
	// //遍历默认桶中的键值对
	// data := db.GetAll("")
	// for k, v := range data {
	// 	if value, ok := v.([][]byte); ok { //数据断言
	// 		fmt.Printf("第%v个数据：, [%s] is %s\n", k, value[0], value[1])
	// 	}

	// }
	// db.Close()

}
