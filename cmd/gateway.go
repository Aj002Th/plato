package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"plato/gateway"
)

func init() {
	rootCmd.AddCommand(gatewayCmd)
}

var gatewayCmd = &cobra.Command{
	Use:   "gateway",
	Short: "plato - gateway 服务",
	Run:   GatewayHandle,
}

func GatewayHandle(cmd *cobra.Command, args []string) {
	fmt.Println("plato gateway")
	gateway.RunMain(ConfigPath)
}
