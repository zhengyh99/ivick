package bc

import (
	"fmt"
	"os"
	"strconv"
)

const (
	//命令说明
	Usage = `
printChain		"正向打印区块链"
getBalance --address DATA	"获取指定地址的余额"
send FROM TO AMOUNT MINER DATA	"由FROM转AMOUNT给TO，由MINER挖矿，同时写入DATA"
newWallet		"获取新钱包"
listWallet		"返回钱包地址列表"
`
)

type CLI struct {
	BC *BlockChain
}

func NewCLI() (cli *CLI) {
	cli = &CLI{
		BC: GetBlockChain("start"),
	}
	return
}

//CLI执行完成后，需及时关闭，释放资源
func (cli *CLI) Close() {
	cli.BC.Close()
}
func (cli *CLI) printChain() {
	bc := cli.BC
	for iter := bc.Iter(); iter.HasNext(); {
		b := iter.Next()
		fmt.Printf("块ID :%x \n", b.Hash)
		for i, txs := range b.TXs {
			fmt.Println("begin================")
			fmt.Printf("第%v个交易，交易ID：%x\n-----------------\n\n", i, txs.TXID)
			fmt.Println("inputs：{")
			for _, in := range txs.TXInputs {
				fmt.Printf("txid:%x,\nindex:%v\n,sig:%x\n\n", in.TXid, in.Index, in.PubKey)
			}
			fmt.Println("}")
			fmt.Println("outputs：{")
			for _, out := range txs.TXOutputs {
				fmt.Printf("pub:%s,\nvalues:%v\n\n", out.PubKeyHash, out.Value)
			}
			fmt.Println("}")
			fmt.Println("end===============")
		}

	}

}

func (cli *CLI) GetBalance(address string) {
	_, total := cli.BC.FindUTXOs(AddrToPubKeyHash(address), -1)
	fmt.Printf("地址：[%s] 的余额为：%.2f", address, total)
}

//send FROM TO AMOUNT MINER DATA	"由FROM转AMOUNT给TO，由MINER挖矿，同时写入DATA"
func (cli *CLI) send(from, to string, amount float64, miner, data string) {
	//fmt.Println("from:", from, ",to:", to, ",amount:", amount, ",miner:", miner, ",data:", data)
	tx := cli.BC.NewTransaction(from, to, amount)
	if tx == nil {
		fmt.Println("无效的交易")
		return
	}
	cbTx := NewCoinBaseTX(miner, data)
	cli.BC.AddBlock([]*Transaction{tx, cbTx})
	fmt.Println("转账成功")
}

func (cli *CLI) NewWallet() {
	ws := NewWallets()
	address, wallet := ws.CreateWallet()
	fmt.Printf(" 公钥 ：%v\n", wallet.PublicKey)
	fmt.Printf(" 私钥 ：%v\n", wallet.PrivateKey)
	fmt.Printf(" 地址 ：%s\n", address)
}
func (cli *CLI) ListWallet() {
	ws := NewWallets()
	fmt.Println("打印地址列表。。。。。")
	for _, addr := range ws.ListAddress() {
		fmt.Printf("address :[%s]\n", addr)

	}
}
func (cli *CLI) Run() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println(Usage)
		return
	}
	switch args[1] {
	case "newWallet":
		cli.NewWallet()
	case "listWallet":
		cli.ListWallet()
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
