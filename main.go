package main

import (
	"blockchain_go/block"
	"blockchain_go/cli"
)

func main() {
	bc := block.NewBlockchain()
	defer bc.DB.Close()

	cli := cli.CLI{Bc: bc}
	cli.Run()
}
