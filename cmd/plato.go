package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

var (
	ConfigPath string
)

func init() {
	cobra.OnInitialize(initConfig)

	// 命令行参数
	rootCmd.PersistentFlags().StringVar(
		&ConfigPath,
		"config",
		"./plato.yaml",
		"config file (default is ./plato.yaml)")
}

func initConfig() {
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Panicln(err)
	}
}

var rootCmd = &cobra.Command{
	Use:   "plato",
	Short: "im system: plato",
	Run:   Plato,
}

func Plato(cmd *cobra.Command, args []string) {
	fmt.Println("plato rootCmd")
}
