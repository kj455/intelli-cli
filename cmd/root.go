package cmd

import (
	"fmt"
	"os"

	"github.com/kj455/intelli-cli/gateway"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "?",
	Short: "IntelliCLI helps you find exact command you want to run",
	Long:  `IntelliCLI to help you find exact command you want to run`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		desc := args[0]

		chat := gateway.CreateChatClient()
		res, err := chat.CreateCompletion(gateway.CreateCompletionRequest{Description: ComposeCLICompletionPrompt(desc)})
		if err != nil {
			fmt.Println("error", err)
			return
		}

		fmt.Println(res.Choices[0].Messages.Content)
	},
}

func ComposeCLICompletionPrompt(desc string) string {
	return `What CLI command will accomplish the following objectives? The answer should follow this format 'command: XXX
note: XXX'. Objectives: ` + desc
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

}
