package cli

import (
	"blockchain_go/internal/blockchain/block"
	"blockchain_go/internal/blockchain/wallet"
	"blockchain_go/pkg/utils"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strconv"
)

var (
	rootCmd = &cobra.Command{
		Use:     "wavepig",
		Long:    "区块链学习！",
		Version: "0.0.1",
	}
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
	rootCmd.AddCommand(
		NewGetAddressCmd(),
		NewPrintCmd(),
		NewCreateWalletCmd(),
		NewCreateBlockCmd(),
		NewListCmd(),
		NewSendCmd(),
		NewReindexUTXOCmd(),
	)
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
			if !wallet.ValidateAddress(address) {
				log.Panic("ERROR: Address is not valid")
			}
			bc := block.NewBlockchain()
			UTXOSet := block.UTXOSet{bc}
			defer bc.DB.Close()

			balance := 0
			pubKeyHash := utils.Base58Decode([]byte(address))
			pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
			UTXOs := UTXOSet.FindUTXO(pubKeyHash)

			for _, out := range UTXOs {
				balance += out.Value
			}

			fmt.Printf("Balance of '%s': %d\n", address, balance)
		},
	}
	return cmd
}

func NewCreateWalletCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create_w",
		Short: "创建钱包",
		Run: func(cmd *cobra.Command, args []string) {
			w := wallet.NewWallet()
			address := w.GetAddress()
			fmt.Printf("Your new address: %s\n", address)

			wallets, err := wallet.NewWallets()
			if err != nil {
				fmt.Println("文件无数据：", err.Error())
				wallets = new(wallet.Wallets)
				wallets.Wallets = make(map[string]*wallet.Wallet)
			}
			wallets.Wallets[string(address)] = w

			wallets.SaveToFile()
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
			defer bc.DB.Close()

			UTXOSet := block.UTXOSet{bc}
			UTXOSet.Reindex()

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
			bc := block.NewBlockchain()
			defer bc.DB.Close()

			bci := bc.Iterator()

			for {
				b := bci.Next()

				fmt.Printf("============ Block %x ============\n", b.Hash)
				fmt.Printf("Prev. block: %x\n", b.PrevBlockHash)
				pow := block.NewProofOfWork(b)
				fmt.Printf("PoW: %s\n\n", strconv.FormatBool(pow.Validate()))
				for _, tx := range b.Transactions {
					fmt.Println(tx)
				}
				fmt.Printf("\n\n")

				if len(b.PrevBlockHash) == 0 {
					break
				}
			}
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

func NewReindexUTXOCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reindex",
		Short: "重建UTXO集",
		Run: func(cmd *cobra.Command, args []string) {
			bc := block.NewBlockchain()
			UTXOSet := block.UTXOSet{bc}
			UTXOSet.Reindex()

			count := UTXOSet.CountTransactions()
			fmt.Printf("Done! There are %d transactions in the UTXO set.\n", count)
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

			bc := block.NewBlockchain()
			UTXOSet := block.UTXOSet{bc}
			defer bc.DB.Close()

			tx := block.NewUTXOTransaction(from, to, int(amount), &UTXOSet)
			cbTx := block.NewCoinbaseTX(from, "")
			txs := []*block.Transaction{cbTx, tx}

			newBlock := bc.MineBlock(txs)
			UTXOSet.Update(newBlock)
			fmt.Println("Success!")
		},
	}
	return cmd
}

func MainCmd() {
	err := rootCmd.Execute()
	must(err)
}
