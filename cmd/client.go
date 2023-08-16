package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"plato/client"
)

func init() {
	rootCmd.AddCommand(clientCmd)
}

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "plato cui 客户端",
	Run:   ClientHandle,
}

func ClientHandle(cmd *cobra.Command, args []string) {
	fmt.Println("plato client")
	client.RunMain()
}
