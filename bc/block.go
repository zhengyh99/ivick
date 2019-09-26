package bc

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"math/big"
	"time"
)

//定义块
type Block struct {
	Version    uint64
	PrivHash   []byte
	MerkelRoot []byte
	TimeStamp  uint64
	Difficulty uint64
	Nonce      uint64
	Hash       []byte
	TXs        []*Transaction
}

//返回sha256 hash
func Sha256Hash(src []byte) [32]byte {
	return sha256.Sum256(src)
}

//创建块
func NewBlock(txs []*Transaction, prevBlockHash []byte) (block *Block) {
	var now time.Time
	block = &Block{
		Version:    0,
		PrivHash:   prevBlockHash,
		TimeStamp:  uint64(now.Unix()),
		Difficulty: 0,
		TXs:        txs,
	}
	block.MakeMerkelRoot()
	pow := NewProofOfWork(block)
	hash, nonce := pow.Run()
	block.Hash = hash
	block.Nonce = nonce
	return
}

//uint64类型转换成[]byte
func (block *Block) uintToBytes(num uint64) []byte {
	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.BigEndian, num)
	if err != nil {
		fmt.Println("binary write error :", err)
	}
	return buffer.Bytes()

}

//计算区块hash值和工作量目标值 ，nonce为难度值
func (block *Block) GetHashAndTarget(nonce uint64) ([]byte, big.Int) {

	tmp := [][]byte{
		block.uintToBytes(block.Version),
		block.PrivHash,
		block.MerkelRoot,
		block.uintToBytes(block.TimeStamp),
		block.uintToBytes(block.Difficulty),
		block.uintToBytes(nonce),
	}
	blockInfo := bytes.Join(tmp, []byte{})
	hash := Sha256Hash(blockInfo)
	tmpInt := big.Int{}
	tmpInt.SetBytes(hash[:])
	return hash[:], tmpInt
}

//计算梅克尔根
func (block *Block) MakeMerkelRoot() {
	var txids []byte
	for _, tx := range block.TXs {
		txids = append(txids, tx.TXID...)
	}
	hash := Sha256Hash(txids)
	block.MerkelRoot = hash[:]
}

//生成创世块
func GenesisBlock(address string) *Block {
	cb := NewCoinBaseTX(address, "创世块证书")
	return NewBlock([]*Transaction{cb}, []byte{})
}

//序列化为[]byte
func (block *Block) Serialize() []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	encoder.Encode(block)
	return buffer.Bytes()
}

//反序列化为Block 实例
func DeSerialize(data []byte) (block Block, err error) {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err = decoder.Decode(&block)
	if err != nil {
		return block, err
	}
	return block, nil
}
