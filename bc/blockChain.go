package bc

import (
	"fmt"
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
func GetBlockChain(address string) *BlockChain {
	db := boltUse.OpenBoltDB(dbFile, bucketName)
	if len(db.GET([]byte("tail"))) == 0 { //数据库中无数据
		genesisBlock := GenesisBlock(address) //将创世块写入数据库
		db.Put(genesisBlock.Hash, genesisBlock.Serialize())
		db.Put([]byte("tail"), genesisBlock.Hash)
	}
	return &BlockChain{
		Blocks: db,
		tail:   db.GET([]byte("tail")),
	}
}

//向区块链中加入区块
func (bc *BlockChain) AddBlock(txs []*Transaction) {
	db := bc.Blocks

	//如果桶不存在，初始化区块链
	block := NewBlock(txs, db.GET([]byte("tail")))
	db.Put(block.Hash, block.Serialize())
	db.Put([]byte("tail"), block.Hash)
	bc.tail = db.GET([]byte("tail"))
}

//根据idHash 返回block
func (bc *BlockChain) GetBlock(idHash []byte) (block Block, err error) {
	db := bc.Blocks
	b := db.GET(idHash)
	block, err = DeSerialize(b)
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

//创建迭代器
func (bc *BlockChain) Iter() BlockChainIter {
	return BlockChainIter{
		BC:          bc,
		CurentBlock: bc.GetTail(),
	}

}

func (bc *BlockChain) FindUTXOs(address string) []TXOutput {
	var UTXO []TXOutput
	spentOutputs := make(map[string][]int64)
	for iter := bc.Iter(); iter.HasNext(); {
		b := iter.Next()
		for _, tx := range b.TXs {
			fmt.Printf("Current txid is %x\n", tx.TXID)
		OUTPUT:
			for i, output := range tx.TXOutputs {
				fmt.Printf("Current index is %v\n", i)
				for _, j := range spentOutputs[string(tx.TXID)] {
					if int64(i) == j { //判断当前output 是否被前一交易的input 所引用（应用代表交易金额已经被花光）
						continue OUTPUT
					}
				}
				if output.PubKeyHash == address {
					UTXO = append(UTXO, output)
				}
			}
			if !tx.IsCoinBase() {
				for _, input := range tx.TXInputs {
					if input.Sig == address {
						indexArray := spentOutputs[string(input.TXid)]
						indexArray = append(indexArray, input.Index)
					}
				}
			}

		}
	}
	return UTXO
}
