package bc

import (
	"zoin/boltUse"
)

const (
	dbFile     = "blockChain.db"
	bucketName = "blockChain"
)

//定义区块链
type BlockChain struct {
	Blocks *boltUse.BoltDB
	tail   []byte
}

func NewBlockChain() *BlockChain {
	db := boltUse.OpenBoltDB(dbFile, bucketName)
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

func (bc *BlockChain) GetBlock(id []byte) (block []byte) {
	db := bc.Blocks
	block = db.GET(id)
	return
}

func (bc *BlockChain) GetTail() (tail []byte) {
	db := bc.Blocks
	tail = db.GET([]byte("tail"))
	return

}

func (bc *BlockChain) GetAll() []interface{} {

	db := bc.Blocks
	return db.GetAll(bucketName)
}

func (bc *BlockChain) Close() {
	bc.Blocks.Close()
}
