package bc

import (
	"zoin/boltUse"
)

const (
	dbFile     = "blockChain.db" //数据库文件
	bucketName = "blockChain"    //数据桶
)

//定义区块链
type BlockChain struct {
	Blocks *boltUse.BoltDB
	tail   []byte
}

//新建区块链
func NewBlockChain() *BlockChain {
	db := boltUse.OpenBoltDB(dbFile, bucketName)
	//判断桶是否存在
	if db.HasBucket(bucketName) == false {
		genesisBlock := GenesisBlock()
		db.Put(genesisBlock.Hash, genesisBlock.Serialize())
		db.Put([]byte("tail"), genesisBlock.Hash)
	}
	return &BlockChain{
		Blocks: db,
		tail:   db.GET([]byte("tail")),
	}
}

//向区块链中加入区块
func (bc *BlockChain) AddBlock(data string) {
	db := bc.Blocks

	//如果桶不存在，初始化区块链
	block := NewBlock(data, db.GET([]byte("tail")))
	db.Put(block.Hash, block.Serialize())
	db.Put([]byte("tail"), block.Hash)
	bc.tail = db.GET([]byte("tail"))
}

//根据idHash 返回block
func (bc *BlockChain) GetBlock(idHash []byte) (block []byte) {
	db := bc.Blocks
	block = db.GET(idHash)
	return
}

//获取最后一个Block的Hash
func (bc *BlockChain) GetTail() (tail []byte) {
	db := bc.Blocks
	tail = db.GET([]byte("tail"))
	return

}

//返回桶中所有数据
func (bc *BlockChain) GetAll() []interface{} {

	db := bc.Blocks
	return db.GetAll(bucketName)
}

//关闭Bolt资源
func (bc *BlockChain) Close() {
	bc.Blocks.Close()
}
