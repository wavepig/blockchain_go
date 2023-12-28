package cli

import (
	"blockchain_go/internal/blockchain/block"
	"blockchain_go/internal/blockchain/wallet"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var (
	rootCmd = &cobra.Command{
		Use:     "wavepig",
		Long:    "区块链学习！",
		Version: "0.0.1",
	}
	cli = CLI{}
)

var (
	address string
	from    string
	to      string
	amount  int64
)

func init() {
	// 默认值设置
	rootCmd.PersistentFlags().StringVarP(&address, "address", "a", "", "地址信息")
	rootCmd.PersistentFlags().StringVarP(&from, "from", "f", "", "转账地址")
	rootCmd.PersistentFlags().StringVarP(&to, "to", "t", "", "目标地址")
	rootCmd.PersistentFlags().Int64VarP(&amount, "amount", "m", 0, "转账金额")
	rootCmd.AddCommand(NewGetAddressCmd())
	rootCmd.AddCommand(NewPrintCmd())
	rootCmd.AddCommand(NewCreateWalletCmd())
	rootCmd.AddCommand(NewCreateBlockCmd())
	rootCmd.AddCommand(NewListCmd())
	rootCmd.AddCommand(NewSendCmd())
}

func must(err error) {
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func NewGetAddressCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "获取地址余额",
		Run: func(cmd *cobra.Command, args []string) {
			if address == "" {
				must(errors.New("需要传递 address"))
			}
			cli.getBalance(address)
		},
	}
	return cmd
}

func NewCreateWalletCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create_w",
		Short: "创建钱包",
		Run: func(cmd *cobra.Command, args []string) {
			//wallets, _ := wallet.NewWallets()
			//address := wallets.CreateWallet()
			//wallets.SaveToFile()
			//
			//fmt.Printf("Your new address: %s\n", address)

			wallet := wallet.NewWallet()
			address := wallet.GetAddress()

			fmt.Printf("Your new address: %s\n", address)
		},
	}
	return cmd
}

func NewCreateBlockCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create_b",
		Short: "创建初始区块",
		Run: func(cmd *cobra.Command, args []string) {
			if address == "" {
				must(errors.New("需要传递 address"))
			}

			if !wallet.ValidateAddress(address) {
				log.Panic("ERROR: Address is not valid")
			}
			bc := block.CreateBlockchain(address)
			bc.DB.Close()
			fmt.Println("Done!")
		},
	}
	return cmd
}

func NewPrintCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "print",
		Short: "获取地址余额",
		Run: func(cmd *cobra.Command, args []string) {
			cli.printChain()
		},
	}
	return cmd
}

func NewListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "列出钱包文件中的所有地址",
		Run: func(cmd *cobra.Command, args []string) {
			wallets, err := wallet.NewWallets()
			if err != nil {
				log.Panic(err)
			}
			addresses := wallets.GetAddresses()

			for _, address := range addresses {
				fmt.Println(address)
			}
		},
	}
	return cmd
}

func NewSendCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send",
		Short: "从地址到地址发送硬币的数量",
		Run: func(cmd *cobra.Command, args []string) {
			if from == "" || to == "" || amount <= 0 {
				log.Panic("ERROR: from and to and amount non null")
			}

			if !wallet.ValidateAddress(from) {
				log.Panic("ERROR: Sender address is not valid")
			}
			if !wallet.ValidateAddress(to) {
				log.Panic("ERROR: Recipient address is not valid")
			}

			bc := block.NewBlockchain(from)
			defer bc.DB.Close()

			tx := block.NewUTXOTransaction(from, to, int(amount), bc)
			bc.MineBlock([]*block.Transaction{tx})
			fmt.Println("Success!")
		},
	}
	return cmd
}

func MainCmd() {
	err := rootCmd.Execute()
	must(err)
}
