package bc

import (
	"crypto/elliptic"
	"fmt"
	"io/ioutil"
)

const walletsFile = "wallets.db"

type Wallets struct {
	WS map[string]*Wallet
}

func NewWallets() (ws *Wallets) {
	ws = &Wallets{
		WS: make(map[string]*Wallet),
	}
	ws.loadFile()
	return
}

func (ws *Wallets) CreateWallet() (address string, wallet *Wallet) {
	wallet = NewWallet()
	address = wallet.NewAddres()
	ws.WS[address] = wallet
	ws.saveToFile()
	return
}

// 根据 地址找到wallet 如果没有返回nil
func (ws *Wallets) FindByAddress(address string) (wallet *Wallet) {
	wallet = ws.WS[address]
	return
}

func (ws *Wallets) saveToFile() {
	RegSerialize(elliptic.P256())
	w := Serialize(ws)
	err := ioutil.WriteFile(walletsFile, w, 0666)
	if err != nil {
		fmt.Println("ioutil write file error:", err)
	}
	fmt.Println("数据已保存")

}

func (ws *Wallets) loadFile() {
	if !PathExists(walletsFile) {
		ws.WS = make(map[string]*Wallet)
		return

	}
	f, err := ioutil.ReadFile(walletsFile)
	if err != nil {
		fmt.Println("ioutil read file error：", err)
	}
	RegSerialize(elliptic.P256())
	var wsTemp Wallets
	err1 := DeSerialize(f, &wsTemp)
	if err1 != nil {
		fmt.Println(" deserialize error：", err1)
	}
	ws.WS = wsTemp.WS
}

func (ws *Wallets) ListAddress() (addresses []string) {

	for key := range ws.WS {
		addresses = append(addresses, key)
	}
	return
}
