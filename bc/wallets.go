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

func (ws *Wallets) CreateWallet() string {
	w := NewWallet()
	address := w.NewAddres()
	ws.WS[address] = w
	ws.saveToFile()
	return address
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
