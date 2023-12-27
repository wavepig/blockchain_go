package cli

import (
	"blockchain_go/internal/blockchain/block"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

// CLI 负责处理命令行参数的CLI
type CLI struct{}

func (cli *CLI) createBlockchain(address string) {
	bc := block.CreateBlockchain(address)
	bc.DB.Close()
	fmt.Println("Done!")
}

func (cli *CLI) getBalance(address string) {
	bc := block.NewBlockchain(address)
	defer bc.DB.Close()

	balance := 0
	UTXOs := bc.FindUTXO(address)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("'%s' 账号余额为: %d\n", address, balance)
}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  get -address ADDRESS - 获取地址的余额")
	fmt.Println("  create -address ADDRESS - 创建一个区块链，并将创世区块奖励发送到ADDRESS")
	fmt.Println("  print - 打印区块链的所有块")
	fmt.Println("  send -from FROM -to TO -amount AMOUNT - 从地址到地址发送硬币的数量")
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) printChain() {
	// TODO: Fix this
	bc := block.NewBlockchain("")
	defer bc.DB.Close()

	bci := bc.Iterator()

	for {
		blockValue := bci.Next()

		fmt.Printf("Prev. hash: %x\n", blockValue.PrevBlockHash)
		fmt.Printf("Hash: %x\n", blockValue.Hash)
		fmt.Printf("Nonce: %d\n", blockValue.Nonce)
		pow := block.NewProofOfWork(blockValue)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		transactions, err := json.MarshalIndent(blockValue.Transactions, "", " ")
		if err == nil {
			fmt.Printf("Data: %v\n", string(transactions))
		} else {
			fmt.Printf("Data: %v\n", err.Error())
		}
		fmt.Println()

		if len(blockValue.PrevBlockHash) == 0 {
			break
		}
	}
}

func (cli *CLI) send(from, to string, amount int) {
	bc := block.NewBlockchain(from)
	defer bc.DB.Close()

	tx := block.NewUTXOTransaction(from, to, amount, bc)
	bc.MineBlock([]*block.Transaction{tx})
	fmt.Println("Success!")
}

// 创建区块 xx create -address [address]
// 打印区块 xx print
// 查看账号余额 xx get -address [address]
// 转账 xx send -from [原地址] -to [目的地址] -amount [多少]
// Run 解析命令行参数并处理命令
func (cli *CLI) Run() {
	cli.validateArgs()

	getBalanceCmd := flag.NewFlagSet("get", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("create", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "获取一个地址的区块集合")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "创建区块")
	sendFrom := sendCmd.String("from", "", "原钱包地址")
	sendTo := sendCmd.String("to", "", "目的钱包地址")
	sendAmount := sendCmd.Int("amount", 0, "转账多少")

	switch os.Args[1] {
	case "get":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "create":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "print":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			os.Exit(1)
		}
		cli.getBalance(*getBalanceAddress)
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			os.Exit(1)
		}
		cli.createBlockchain(*createBlockchainAddress)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			os.Exit(1)
		}

		cli.send(*sendFrom, *sendTo, *sendAmount)
	}
}
