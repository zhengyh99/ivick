package bc

import (
	"bytes"
	"encoding/gob"
)

const reward = 12.5

//定义交易结构
type Transaction struct {
	TXID      []byte     //交易ID
	TXInputs  []TXInput  //交易输入数组
	TXOutputs []TXOutput //交易输出数组
}

//定义交易输入
type TXInput struct {
	TXid  []byte //引用所在交易ID
	Index int64  //引用output的索引值
	Sig   string //解锁脚本
}

//定义交易输出
type TXOutput struct {
	value      float64 //转账金额
	PubKeyHash string  //锁定脚本
}

//Hash计算交易ID
func (tx *Transaction) SetHash() {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(tx)
	if err != nil {
		panic(err)
	}
	data := Sha256Hash(buffer.Bytes())
	tx.TXID = data[:]
}

//提供挖矿交易方法
func NewCoinBaseTX(address, sig string) (tx *Transaction) {
	//挖矿交易的特点：
	//1、只有一个input 和一个output
	//2、无需引用交易id
	//3、无需引用index
	//4、无需指定签名，sig可同矿工自定义，一般填写矿工池的名称
	input := TXInput{TXid: []byte{}, Index: -1, Sig: sig}
	output := TXOutput{value: reward, PubKeyHash: address}
	tx = &Transaction{TXInputs: []TXInput{input}, TXOutputs: []TXOutput{output}}
	tx.SetHash()
	return

}
