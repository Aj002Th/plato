package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"plato/ipconfig"
)

func init() {
	rootCmd.AddCommand(ipConfigCmd)
}

var ipConfigCmd = &cobra.Command{
	Use:   "ipconfig",
	Short: "plato - ipconfig 服务",
	Run:   IpConfigHandle,
}

func IpConfigHandle(cmd *cobra.Command, args []string) {
	fmt.Println("plato ipconfig")
	ipconfig.RunMain(ConfigPath)
}
