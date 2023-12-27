package cli

import (
	"blockchain_go/internal/blockchain/block"
	"blockchain_go/internal/blockchain/wallet"
	"blockchain_go/pkg/utils"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

// CLI 负责处理命令行参数的CLI
type CLI struct{}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  get -address ADDRESS - 获取地址的余额")
	fmt.Println("  create_b -address ADDRESS - 创建一个区块链，并将创世区块奖励发送到ADDRESS")
	fmt.Println("  create_w - 生成一个新的密钥对并将其保存到钱包文件中")
	fmt.Println("  list - 列出钱包文件中的所有地址")
	fmt.Println("  print - 打印区块链的所有块")
	fmt.Println("  send -from FROM -to TO -amount AMOUNT - 从地址到地址发送硬币的数量")
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) printChainA() {
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

// 创建区块 xx create -address [address]
// 打印区块 xx print
// 查看账号余额 xx get -address [address]
// 转账 xx send -from [原地址] -to [目的地址] -amount [多少]
// Run 解析命令行参数并处理命令
func (cli *CLI) Run() {
	cli.validateArgs()

	getBalanceCmd := flag.NewFlagSet("get", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("create_b", flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("create_w", flag.ExitOnError)
	listAddressesCmd := flag.NewFlagSet("list", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "余额获取")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "创建区块")
	sendFrom := sendCmd.String("from", "", "本地地址")
	sendTo := sendCmd.String("to", "", "目的地址")
	sendAmount := sendCmd.Int("amount", 0, "转账数量")

	switch os.Args[1] {
	case "get":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "create_b":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "create_w":
		err := createWalletCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "list":
		err := listAddressesCmd.Parse(os.Args[2:])
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

	if createWalletCmd.Parsed() {
		cli.createWallet()
	}

	if listAddressesCmd.Parsed() {
		cli.listAddresses()
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

// 创建账号
func (cli *CLI) createBlockchain(address string) {
	if !wallet.ValidateAddress(address) {
		log.Panic("ERROR: Address is not valid")
	}
	bc := block.CreateBlockchain(address)
	bc.DB.Close()
	fmt.Println("Done!")
}

// 创建钱包
func (cli *CLI) createWallet() {
	wallets, _ := wallet.NewWallets()
	address := wallets.CreateWallet()
	wallets.SaveToFile()

	fmt.Printf("Your new address: %s\n", address)
}

// 获取余额
func (cli *CLI) getBalance(address string) {
	if !wallet.ValidateAddress(address) {
		log.Panic("ERROR: Address is not valid")
	}
	bc := block.NewBlockchain(address)
	defer bc.DB.Close()

	balance := 0
	pubKeyHash := utils.Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	UTXOs := bc.FindUTXO(pubKeyHash)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of '%s': %d\n", address, balance)
}

func (cli *CLI) listAddresses() {
	wallets, err := wallet.NewWallets()
	if err != nil {
		log.Panic(err)
	}
	addresses := wallets.GetAddresses()

	for _, address := range addresses {
		fmt.Println(address)
	}
}

func (cli *CLI) printChain() {
	bc := block.NewBlockchain("")
	defer bc.DB.Close()

	bci := bc.Iterator()

	for {
		blockHash := bci.Next()

		fmt.Printf("============ Block %x ============\n", blockHash.Hash)
		fmt.Printf("Prev. hash: %x\n", blockHash.PrevBlockHash)
		fmt.Printf("Hash: %x\n", blockHash.Hash)
		fmt.Printf("Nonce: %d\n", blockHash.Nonce)
		pow := block.NewProofOfWork(blockHash)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println("=============  end  ======================")
		//transactions, err := json.MarshalIndent(blockHash.Transactions, "", " ")
		//if err == nil {
		//	fmt.Printf("Data: %v\n", string(transactions))
		//} else {
		//	fmt.Printf("Data: %v\n", err.Error())
		//}
		fmt.Printf("\n\n")

		if len(blockHash.PrevBlockHash) == 0 {
			break
		}
	}
}

func (cli *CLI) send(from, to string, amount int) {
	if !wallet.ValidateAddress(from) {
		log.Panic("ERROR: Sender address is not valid")
	}
	if !wallet.ValidateAddress(to) {
		log.Panic("ERROR: Recipient address is not valid")
	}

	bc := block.NewBlockchain(from)
	defer bc.DB.Close()

	tx := block.NewUTXOTransaction(from, to, amount, bc)
	bc.MineBlock([]*block.Transaction{tx})
	fmt.Println("Success!")
}
