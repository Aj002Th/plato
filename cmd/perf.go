package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"plato/perf"
)

func init() {
	rootCmd.AddCommand(perfCmd)
	perfCmd.PersistentFlags().Int32Var(&perf.TcpConnNum, "tcp_conn_num", 10000, "tcp 连接的数量，默认10000")
}

var perfCmd = &cobra.Command{
	Use:   "perf",
	Short: "plato -  perf 连接压测",
	Run:   PerfHandle,
}

func PerfHandle(cmd *cobra.Command, args []string) {
	fmt.Println("plato perf")
	perf.RunMain()
}
