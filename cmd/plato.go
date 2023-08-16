package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

func init() {
	cobra.OnInitialize(initConfig)
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
