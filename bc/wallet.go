package bc

//git clone https://github.com/golang/crypto.git
//git clone https://github.com/btcsuite/btcutil.git
import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"z/btcutil/base58"
)

type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  []byte
}

func NewWallet() *Wallet {
	curve := elliptic.P256()
	privKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		fmt.Println("ecdsa generatekey error:", err)
	}
	pubkey := append(privKey.PublicKey.X.Bytes(), privKey.PublicKey.Y.Bytes()...)
	return &Wallet{
		PrivateKey: privKey,
		PublicKey:  pubkey,
	}

}

func (w *Wallet) NewAddres() (address string) {
	pubKey := w.PublicKey
	hash := Sha256Hash(pubKey)
	//返回 ripemd160 编码
	rip160HashValue := Rip160Hash(hash[:])
	version := byte(00)
	payLoad := append([]byte{version}, rip160HashValue...)
	hash1 := Sha256Hash(payLoad)
	hash2 := Sha256Hash(hash1[:])
	checkCode := hash2[:4]
	payLoad = append(payLoad, checkCode...)
	address = base58.Encode(payLoad)
	return
}
