package bc

import (
	"math/big"
)

//定义工作量
type ProofOfWork struct {
	block  *Block
	target *big.Int
}

//为区块中加新的工作量
func NewProofOfWork(block *Block) ProofOfWork {
	pow := ProofOfWork{
		block: block,
	}
	targetStr := "0000100000000000000000000000000000000000000000000000000000000000"
	tmpInt := big.Int{}
	tmpInt.SetString(targetStr, 16)
	pow.target = &tmpInt
	return pow
}

//工作量运算，返回区块hash和难度值
func (pow *ProofOfWork) Run() (hash []byte, nonce uint64) {
	var tmpTarget big.Int
	for {
		hash, tmpTarget = pow.block.GetHashAndTarget(nonce)

		if tmpTarget.Cmp(pow.target) == -1 {
			//fmt.Printf("Nonce is reserch ,Hash :%v,Nonce :%v\n", hash, nonce)
			return
		} else {
			nonce++
		}
	}

}
