package bc

import (
	"bytes"
	"crypto/ecdsa"
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
	//交易签名检验
	for _, tx := range txs {
		// fmt.Println(tx)
		// fmt.Println("====\n====")
		if !bc.VerifyTransaction(tx) {
			fmt.Println("矿工发现交易验证失败！！！")
			return
		}
	}
	db := bc.Blocks

	//如果桶不存在，初始化区块链
	block := NewBlock(txs, db.GET([]byte("tail")))
	//写入数据库
	db.Put(block.Hash, block.Serialize())
	db.Put([]byte("tail"), block.Hash)
	bc.tail = db.GET([]byte("tail")) //读取区块链尾
}

//根据idHash 返回block
func (bc *BlockChain) GetBlock(idHash []byte) (block *Block, err error) {
	db := bc.Blocks
	b := db.GET(idHash)
	err1 := DeSerialize(b, &block)
	if err1 != nil {
		err = err1
		return
	}

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
func (bc *BlockChain) Iter() *BlockChainIter {
	return &BlockChainIter{
		BC:          bc,
		CurentBlock: bc.GetTail(),
	}

}

//查找指定金额：needAmount的UTXO（若needAmount==-1:代表查找全部UTXO）
//返回找到UTXO所在交易ID 和index切片 组成map 及 所得的金额总额
func (bc *BlockChain) FindUTXOs(pubKeyHash []byte, needAmount float64) (utxos map[string][]int64, balance float64) {
	spentOutputs := make(map[string][]int64)
	utxos = make(map[string][]int64)
SUFFICIENT:
	for iter := bc.Iter(); iter.HasNext(); {
		b := iter.Next()

		for _, tx := range b.TXs {
			// fmt.Printf("Current txid is %x\n", tx.TXID)
		OUTPUT:
			for i, output := range tx.TXOutputs {
				// fmt.Printf("Current index is %v\n", i)
				for _, j := range spentOutputs[string(tx.TXID)] {
					if int64(i) == j { //判断当前output 是否被前一交易的input 所引用（应用代表交易金额已经被花光）
						continue OUTPUT
					}
				}
				if bytes.Equal(output.PubKeyHash, pubKeyHash) {
					balance = balance + output.Value
					utxos[string(tx.TXID)] = append(utxos[string(tx.TXID)], int64(i))
					//当金额充足时停止查找，退出整个循环
					if balance >= needAmount && needAmount != -1 {
						break SUFFICIENT
					}
				}
			}

			//挖矿交易，不执行
			if !tx.IsCoinBase() {

				for _, input := range tx.TXInputs {

					if bytes.Equal(HashPubKey(input.PubKey), pubKeyHash) {
						spentOutputs[string(input.TXid)] = append(spentOutputs[string(input.TXid)], input.Index)
					}
				}
			}

		}
	}
	return

}

//新建普通交易 needAmount 交易成交所需要的金额
func (bc *BlockChain) NewTransaction(from, to string, needAmount float64) (tx *Transaction) {
	ws := NewWallets()
	if ws == nil {
		fmt.Println("地址不存在")
		return nil
	}
	wallet := ws.FindByAddress(from)
	//查找可用金额(即余额) balance
	utxos, balance := bc.FindUTXOs(HashPubKey(wallet.PublicKey), needAmount)
	if balance < needAmount {
		fmt.Println("金额不足，交易失败")
		return nil
	}
	var inputs []TXInput
	var outputs []TXOutput
	//创建 TXInput
	for txid, indexs := range utxos {
		for _, i := range indexs {
			input := TXInput{TXid: []byte(txid), Index: i, Signature: nil, PubKey: wallet.PublicKey}
			inputs = append(inputs, input)
		}
	}
	//创建 TXOutput
	output := NewTXOutput(to, needAmount)
	outputs = append(outputs, *output)
	//找零交易,播放TXOutput
	if balance > needAmount {
		outputs = append(outputs, TXOutput{Value: balance - needAmount, PubKeyHash: HashPubKey(wallet.PublicKey)})
	}
	tx = &Transaction{
		TXInputs:  inputs,
		TXOutputs: outputs,
	}
	//生成TXID
	tx.SetHash()
	//生成交易签名
	bc.SignTransaction(tx, wallet.PrivateKey)
	return
}

//根据TXID查找Transaction
func (bc *BlockChain) FindTransactionByTXID(txid []byte) *Transaction {
	for iter := bc.Iter(); iter.HasNext(); {
		block := iter.Next()
		for _, tx := range block.TXs {
			if bytes.Equal(tx.TXID, txid) {
				return tx
			}
		}
	}
	//未找到
	return nil
}

//生成交易签名
func (bc *BlockChain) SignTransaction(tx *Transaction, privetKey *ecdsa.PrivateKey) {
	prevTXs := make(map[string]Transaction)
	//找出交易中所有的input交易中TXID，对应该的交易，填入map
	for _, input := range tx.TXInputs {
		t := bc.FindTransactionByTXID(input.TXid)
		if t != nil {
			prevTXs[string(input.TXid)] = *t
		}
	}
	//生成签名
	tx.Sign(privetKey, prevTXs)

}

//检验交易签名
func (bc *BlockChain) VerifyTransaction(tx *Transaction) bool {
	//挖矿交易直接路过校验返回true
	if tx.IsCoinBase() {
		return true
	}
	prevTXs := make(map[string]Transaction)
	//找出交易中所有的input交易中TXID，对应该的交易，填入map
	for _, input := range tx.TXInputs {
		t := bc.FindTransactionByTXID(input.TXid)
		if t != nil {
			prevTXs[string(input.TXid)] = *t
		}
	}
	//返回检验结果
	return tx.Verify(prevTXs)
}
