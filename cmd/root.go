package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "contract-cloner",
	Short: "A CLI tool to download verified smart contract source code.",
	Long:  `A command-line tool written in Go to download verified source code for smart contracts from various block explorers.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
