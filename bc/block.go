package bc

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
)

type Block struct {
	Version    uint64
	PrivHash   []byte
	MerkelRoot []byte
	TimeStamp  uint64
	Difficulty uint64
	Hash       []byte
	Data       []byte
}

func (block *Block) Uint64toBytes(num uint64) []byte {
	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.BigEndian, num)
	if err != nil {
		fmt.Println("binary write error:", err)
	}
	return buffer.Bytes()

}
func NewBlock(data string, prevBlockHash []byte) (block *Block) {
	block = &Block{
		Version:    0,
		PrivHash:   prevBlockHash,
		MerkelRoot: []byte{},
		TimeStamp:  0,
		Difficulty: 0,
		Hash:       []byte{},
		Data:       []byte(data),
	}
	block.setHash()
	return
}

func (block *Block) setHash() {

	tmp := [][]byte{
		block.Uint64toBytes(block.Version),
		block.PrivHash,
		block.MerkelRoot,
		block.Uint64toBytes(block.TimeStamp),
		block.Uint64toBytes(block.Difficulty),
		block.Hash,
		block.Data,
	}
	blockInfo := bytes.Join(tmp, []byte{})
	hash := sha256.Sum256(blockInfo)
	block.Hash = hash[:]
}
func GenesisBlock() *Block {
	return NewBlock("创世块", []byte{})
}
