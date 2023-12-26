package main

import (
	"blockchain_go/internal/blockchain/block"
	"blockchain_go/internal/blockchain/cli"
)

func main() {
	bc := block.NewBlockchain()
	defer bc.DB.Close()

	cli := cli.CLI{Bc: bc}
	cli.Run()
}
