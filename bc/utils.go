package bc

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"os"
	"z/btcutil/base58"

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

//地址转public key hash
func AddrToPubKeyHash(address string) (pubKeyHash []byte) {
	addrBytes := base58.Decode(address)
	len := len(addrBytes)
	pubKeyHash = addrBytes[1 : len-4]
	return
}

//将数据进行二次hash 返回前四个byte
func DoubleHash4(data []byte) []byte {
	hash1 := Sha256Hash(data)
	hash2 := Sha256Hash(hash1[:])
	return hash2[:4]
}

//验证地址的有效性
func IsValidAddress(address string) bool {
	addrBytes := base58.Decode(address)
	if len(addrBytes) < 4 {
		return false
	}
	payLoad := addrBytes[:len(addrBytes)-4]
	checkCode1 := addrBytes[len(addrBytes)-4:]
	checkCode2 := DoubleHash4(payLoad)
	return bytes.Equal(checkCode1, checkCode2)

}
