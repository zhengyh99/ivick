package main

import (
	"zoin/bc"
)

func main() {
	//新建 CLI
	cli := bc.NewCLI()
	//程序运行结束后，及时关闭资源
	defer cli.Close()
	//运行CLI
	cli.Run()
}
