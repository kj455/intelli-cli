package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "intelli-cli",
	Short: "IntelliCLI helps you find exact command you want to run",
	Long:  `IntelliCLI to help you find exact command you want to run`,
}

func init() {
	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(questionCmd)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
