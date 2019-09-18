package boltUse

import (
	"bolt"
	"fmt"
)

//定义结构体
type BoltDB struct {
	db     *bolt.DB //数据库
	bucket []byte   //桶
}

//打开数据库文件
func OpenBoltDB(fileName string) *BoltDB {
	boltdb, err := bolt.Open(fileName, 0600, nil)
	if err != nil {
		fmt.Println("bolt open error:", err)
	}
	MyDB := &BoltDB{
		db: boltdb,
	}
	MyDB.SetBucket("") //初始化桶
	return MyDB
}

//关闭数据库
func (this *BoltDB) Close() {
	this.db.Close()
}

//设置当前桶，有使用，无创建
func (this *BoltDB) SetBucket(bucketName string) error {
	if bucketName == "" {
		bucketName = "default"
	}
	bucket := []byte(bucketName)
	if err := this.db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(bucket); err != nil { //判断桶是否存在，不存在建新
			fmt.Println("create failed", err.Error())
			return err
		}
		return nil
	}); err != nil {
		fmt.Println("update error is:", err)
	}
	this.bucket = bucket
	return nil

}

//判断桶是否存在
func (this *BoltDB) HasBucket(bucketName string) bool {
	if bucketName == "" {
		bucketName = "default"
	}
	bucket := []byte(bucketName)
	has := false
	fmt.Printf("bucket:%s\n", bucket)
	this.db.Update(func(tx *bolt.Tx) error {
		if b := tx.Bucket(bucket); b != nil {

			has = true
		}
		return nil
	})
	return has
}

//添加/修改新的键值
func (this *BoltDB) Put(key, value []byte) error {
	if err := this.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(this.bucket)
		err := b.Put(key, value) //添加 健值
		return err
	}); err != nil {
		return err
	}
	return nil
}

//依键得值
func (this *BoltDB) GET(key []byte) (value []byte) {
	if err := this.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(this.bucket)
		value = b.Get(key) //依健查值
		return nil
	}); err != nil {
		fmt.Println("view error :", err.Error())
		return []byte("")
	}
	return
}

//查询指定桶中的所有键值
func (this *BoltDB) GetAll(bucketName string) (datas []interface{}) {
	fmt.Println("-----------------------------------------")
	this.SetBucket(bucketName)
	this.db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		fmt.Printf("this bucket:%s\n", this.bucket)
		b := tx.Bucket(this.bucket)
		b.ForEach(func(k, v []byte) error { //遍历
			fmt.Println("k:", k, ";v:", v)
			kv := [][]byte{k, v}
			datas = append(datas, kv)
			return nil
		})
		return nil
	})
	return
}
