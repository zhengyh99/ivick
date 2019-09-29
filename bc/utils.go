package bc

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"os"

	"golang.org/x/crypto/ripemd160"
)

//返回sha256 hash
func Sha256Hash(src []byte) [32]byte {
	return sha256.Sum256(src)
}

//返回 ripemd160 编码
func Rip160Hash(src []byte) (result []byte) {
	hasher := ripemd160.New()
	_, err := hasher.Write(src)
	if err != nil {
		fmt.Println("ripemd 160 write error:", err)
	}
	result = hasher.Sum(nil)
	return
}

//序列化interface{} 为[]byte
func Serialize(e interface{}) []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	encoder.Encode(e)
	return buffer.Bytes()
}

//反序列化为Block 实例
func DeSerialize(data []byte, e interface{}) (err error) {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err = decoder.Decode(e)
	if err != nil {
		return
	}
	return
}

//gob编码注册被编码数据不的interface
func RegSerialize(value interface{}) {
	gob.Register(value)
}

//判断文件或目录 文件路径是否存在
func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}
