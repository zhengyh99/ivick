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
	Data       []byte
}

//创建块
func NewBlock(data string, prevBlockHash []byte) (block *Block) {
	var now time.Time
	block = &Block{
		Version:    0,
		PrivHash:   prevBlockHash,
		MerkelRoot: []byte{},
		TimeStamp:  uint64(now.Unix()),
		Difficulty: 0,
		Data:       []byte(data),
	}
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
		block.Hash,
		block.Data,
	}
	blockInfo := bytes.Join(tmp, []byte{})
	hash := sha256.Sum256(blockInfo)
	tmpInt := big.Int{}
	tmpInt.SetBytes(hash[:])
	return hash[:], tmpInt
}

//生成创世块
func GenesisBlock() *Block {
	hash := sha256.Sum256([]byte("创世块"))
	return NewBlock("创世块", hash[:])
}

//序列化为[]byte
func (block *Block) Serialize() []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	encoder.Encode(block)
	return buffer.Bytes()
}

//反序列化为Block 实例
func DeSerliaze(data []byte) (block Block, err error) {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err = decoder.Decode(&block)
	if err != nil {
		return block, err
	}
	return block, nil
}
