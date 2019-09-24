package bc

import (
	"fmt"
	"os"
	"strconv"
)

const (
	//命令说明
	Usage = `
addBlock --data DATA	"添加区块"
printChain		"正向打印区块链"
getBalance --address DATA	"获取指定地址的余额"
send FROM TO AMOUNT MINER DATA	"由FROM转AMOUNT给TO，由MINER挖矿，同时写入DATA"
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

func (cli *CLI) printChain() {
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

//send FROM TO AMOUNT MINER DATA	"由FROM转AMOUNT给TO，由MINER挖矿，同时写入DATA"
func (cli *CLI) send(from, to string, amount float64, miner, data string) {
	tx := cli.BC.NewTransaction(from, to, amount)
	if tx == nil {
		fmt.Println("无效的交易")
		return
	}
	cbTx := NewCoinBaseTX(miner, data)
	cli.BC.AddBlock([]*Transaction{tx, cbTx})
	fmt.Println("转账成功")
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
	case "printChain":
		fmt.Println("======显示区块链中的数据===========")
		cli.printChain()
	case "send":
		if len(args) != 7 {
			fmt.Println("参数个数错误。")
			fmt.Println(Usage)
			return
		}
		//send FROM TO AMOUNT MINER DATA	"由FROM转AMOUNT给TO，由MINER挖矿，同时写入DATA"
		from := args[2]
		to := args[3]
		amount, err := strconv.ParseFloat(args[4], 64)
		if err != nil {
			fmt.Println("参数 AMOUNT错误：整数或小数")
			fmt.Println(Usage)
		}
		miner := args[5]
		data := args[6]
		cli.send(from, to, amount, miner, data)
	default:
		fmt.Println("无效命令，请检查")
		fmt.Println(Usage)

	}

}
