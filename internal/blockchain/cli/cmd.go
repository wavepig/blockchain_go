package cli

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
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

func MainCmd() {
	err := rootCmd.Execute()
	must(err)
}
