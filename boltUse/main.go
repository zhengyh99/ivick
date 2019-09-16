package main

import (
	"bolt"
	"fmt"
)

func main() {
	db, err := bolt.Open("mydb.db", 0600, nil)
	if err != nil {
		fmt.Println("bolt open error:", err)
	}
	defer db.Close()
	// if err := db.Update(func(tx *bolt.Tx) error {
	// 	if _, err := tx.CreateBucketIfNotExists([]byte("kl")); err != nil { //判断是否存在
	// 		fmt.Println("create failed", err.Error())
	// 		return err
	// 	}
	// 	b := tx.Bucket([]byte("kl"))
	// 	err = b.Put([]byte("konglong"), []byte("恐龙"))
	// 	b.Put([]byte("konglong2"), []byte("恐龙2"))
	// 	return err
	// }); err != nil {
	// 	fmt.Println("update error is:", err)
	// }

	// if err := db.View(func(tx *bolt.Tx) error {
	// 	b := tx.Bucket([]byte("kl"))
	// 	v := b.Get([]byte("konglong"))
	// 	fmt.Printf("the data is :%s\n", v)
	// 	return nil
	// }); err != nil {
	// 	fmt.Println("view error :", err.Error())
	// }
	// if err := db.View(func(tx *bolt.Tx) error {
	// 	b := tx.Bucket([]byte("kl"))
	// 	v := b.Get([]byte("konglong2"))
	// 	fmt.Printf("the data is :%s\n", v)
	// 	return nil
	// }); err != nil {
	// 	fmt.Println("view error:", err.Error())
	// }

	db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("kl"))
		b.ForEach(func(k, v []byte) error {
			fmt.Printf("key=%s, value=%s\n", k, v)
			return nil
		})
		return nil
	})

}
