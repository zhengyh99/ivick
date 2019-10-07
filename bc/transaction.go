package bc

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/gob"
	"fmt"
	"math/big"
	"strings"
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
func (output *TXOutput) Lock(address string) {

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

//挖矿交易
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

//交易拷贝
func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TXInput
	var outputs []TXOutput
	for _, input := range tx.TXInputs {
		inputs = append(inputs, TXInput{TXid: input.TXid, Index: input.Index, Signature: nil, PubKey: nil})
	}
	for _, output := range tx.TXOutputs {
		outputs = append(outputs, output)
	}
	return Transaction{TXID: tx.TXID, TXInputs: inputs, TXOutputs: outputs}

}

//对交易进行签名
func (tx *Transaction) Sign(privateKey *ecdsa.PrivateKey, prevTXs map[string]Transaction) {
	//如果交易为挖矿交易，不需要签名
	if tx.IsCoinBase() {
		return
	}
	txCopy := tx.TrimmedCopy()
	for i, input := range tx.TXInputs {
		//返回input.TXid对应该的交易
		prevTX := prevTXs[string(input.TXid)]
		//将交易TXOutput的pubkehyHash 赋值给拷贝交易的TXInput的pubkey
		txCopy.TXInputs[i].PubKey = prevTX.TXOutputs[input.Index].PubKeyHash
		//生成拷贝交易的TXID
		txCopy.SetHash()
		txCopy.TXInputs[i].PubKey = nil
		//生成ecdsa签名
		r, s, err := ecdsa.Sign(rand.Reader, privateKey, txCopy.TXID)
		if err != nil {
			panic(err)
		}
		//将签名写入TX TXInput的Signatrue中
		tx.TXInputs[i].Signature = append(r.Bytes(), s.Bytes()...)
	}

}

//交易签名验证
func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	//挖矿交易不需要验证
	if tx.IsCoinBase() {
		return true
	}
	//生成交易拷贝
	txCopy := tx.TrimmedCopy()
	for i, input := range tx.TXInputs {
		prevTX := prevTXs[string(input.TXid)]
		//将交易TXOutput的pubkehyHash 赋值给拷贝交易的TXInput的pubkey
		txCopy.TXInputs[i].PubKey = prevTX.TXOutputs[input.Index].PubKeyHash
		//生成拷贝交易的TXID
		txCopy.SetHash()

		//声明变量
		var r, s, x, y big.Int

		//获取并解析签名
		signature := input.Signature

		r.SetBytes(signature[0 : len(signature)/2])
		s.SetBytes(signature[len(signature)/2:])

		//获取并解析公钥
		pubkey := input.PubKey
		x.SetBytes(pubkey[0 : len(pubkey)/2])
		y.SetBytes(pubkey[len(pubkey)/2:])
		//生成ecdsa公钥
		oPubKey := &ecdsa.PublicKey{
			Curve: elliptic.P256(),
			X:     &x,
			Y:     &y,
		}
		//校验签名
		if !ecdsa.Verify(oPubKey, txCopy.TXID, &r, &s) {
			return false
		}
	}
	return true
}

//格式化Transaction
func (tx *Transaction) String() string {
	var lines []string
	lines = append(lines, fmt.Sprintf("----Transaction ID :%x", tx.TXID))
	for i, input := range tx.TXInputs {
		lines = append(lines, fmt.Sprintf("		Input No. %d", i))
		lines = append(lines, fmt.Sprintf("		TXid :%x", input.TXid))
		lines = append(lines, fmt.Sprintf("		Index  :%d", input.Index))
		lines = append(lines, fmt.Sprintf("		Signature :%x", input.Signature))
		lines = append(lines, fmt.Sprintf("		PubKey :%x", input.PubKey))
	}

	for i, output := range tx.TXOutputs {
		lines = append(lines, fmt.Sprintf("		Output No. %d", i))
		lines = append(lines, fmt.Sprintf("		Value :%f", output.Value))
		lines = append(lines, fmt.Sprintf("		publicKeyHash  :%x", output.PubKeyHash))

	}

	return strings.Join(lines, "\n")
}
