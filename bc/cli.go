package bc

import (
	"fmt"
	"os"
)

const (
	//命令说明
	Usage = `
addBlock --data DATA	"add a block to block chain"
printData		"print all block chain data"
getBalance --address DATA	"get balance by address"
	`
)

type CLI struct {
	BC *BlockChain
}

func NewCLI() (cli *CLI) {
	cli = &CLI{
		BC: GetBlockChain("创世块"),
	}
	return
}

func (cli *CLI) PrintDATA() {
	bc := cli.BC
	for iter := bc.Iter(); iter.HasNext(); {
		b := iter.Next()
		fmt.Printf("b.PrivHash：%x\n", b.PrivHash)
		fmt.Printf("b.Hash:%x \n", b.Hash)
		fmt.Printf("b.daga:%s\n", b.TXs[0].TXInputs[0].Sig)
		fmt.Println("==========")
	}
	bc.Close()
}

func (cli *CLI) GetBalance(address string) {
	utxo := cli.BC.FindUTXOs(address)
	total := 0.0
	for _, v := range utxo {
		total = total + v.Value
	}
	fmt.Printf("地址：[%s] 的余额为：%.2f", address, total)
}

func (cli *CLI) Run() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println(Usage)
		return
	}
	switch args[1] {
	case "addBlock":
		if len(args) == 4 && args[2] == "--data" {
			//cli.BC.AddBlock(args[3])
			fmt.Println("块数据添加成功")
		}
	case "getBalance":
		if len(args) == 4 && args[2] == "--address" {
			cli.GetBalance(args[3])
		}
	case "printData":
		fmt.Println("======显示区块链中的数据===========")
		cli.PrintDATA()
	default:
		fmt.Println("无效命令，请检查")
		fmt.Println(Usage)

	}

}
