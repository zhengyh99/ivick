package bc

import (
	"bytes"
	"encoding/gob"
)

const reward = 12.5

//定义交易块结构
type Transaction struct {
	TXID      []byte     //交易ID
	TXInputs  []TXInput  //交易输入数组
	TXOutputs []TXOutput //交易输出数组
}

//定义交易输入
type TXInput struct {
	TXid      []byte //引用所在交易ID
	Index     int64  //引用output的索引值
	Signature []byte //真正的数字签名 由r s拼成的[]byte
	PubKey    []byte //公钥
}

//定义交易输出
type TXOutput struct {
	Value      float64 //转账金额
	PubKeyHash []byte  //收款方的公钥HASH
}

func NewTXOutput(address string, value float64) (output *TXOutput) {
	output = &TXOutput{
		Value: value,
	}
	output.Lock(address)
	return
}

//锁定公钥HASH
func (output TXOutput) Lock(address string) {

	output.PubKeyHash = AddrToPubKeyHash(address)

}

//Hash计算交易ID
func (tx *Transaction) SetHash() {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(tx)
	if err != nil {
		panic(err)
	}
	//将序列化的Transaction对象Hash
	data := Sha256Hash(buffer.Bytes())
	tx.TXID = data[:]
}

//提供挖矿交易方法
func NewCoinBaseTX(address, data string) (tx *Transaction) {
	//挖矿交易的特点：
	//1、只有一个input 和一个output
	//2、无需引用交易id
	//3、无需引用index
	//4、无需指定签名，sig可同矿工自定义，一般填写矿工池的名称
	input := TXInput{TXid: []byte{}, Index: -1, Signature: nil, PubKey: []byte(data)}
	output := NewTXOutput(address, reward)
	tx = &Transaction{TXInputs: []TXInput{input}, TXOutputs: []TXOutput{*output}}
	tx.SetHash()
	return

}

//判断是否为挖矿交易
func (tx *Transaction) IsCoinBase() bool {
	//挖矿交易的特点：
	//1、只有一个input 和一个output
	//2、无需引用交易id
	//3、无需引用index

	if len(tx.TXInputs) == 1 {
		input := tx.TXInputs[0]
		if bytes.Equal(input.TXid, []byte{}) && input.Index == -1 {
			return true
		}
	}
	return false
}
