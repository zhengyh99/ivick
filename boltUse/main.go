package main

import (
	"bolt"
	"fmt"
)

func main() {
	//打开数据文件
	db, err := bolt.Open("mydb.db", 0600, nil)
	if err != nil {
		fmt.Println("bolt open error:", err)
	}
	defer db.Close()

	//添加键值
	if err := db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte("kl")); err != nil { //判断桶是否存在，不存在建新
			fmt.Println("create failed", err.Error())
			return err
		}
		b := tx.Bucket([]byte("kl"))
		err = b.Put([]byte("konglong"), []byte("恐龙")) //添加 健值
		b.Put([]byte("konglong2"), []byte("恐龙2"))
		return err
	}); err != nil {
		fmt.Println("update error is:", err)
	}

	//查询键-值

	if err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("kl"))
		v := b.Get([]byte("konglong")) //依健查值
		fmt.Printf("the data is :%s\n", v)
		return nil
	}); err != nil {
		fmt.Println("view error :", err.Error())
	}
	if err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("kl"))
		v := b.Get([]byte("konglong2"))
		fmt.Printf("the data is :%s\n", v)
		return nil
	}); err != nil {
		fmt.Println("view error:", err.Error())
	}

	//遍历指定桶中的所有键值
	db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("kl"))
		b.ForEach(func(k, v []byte) error { //遍历
			fmt.Printf("key=%s, value=%s\n", k, v)
			return nil
		})
		return nil
	})

}
