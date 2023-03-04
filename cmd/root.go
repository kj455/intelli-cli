package cmd

import (
	"fmt"
	"os"
	"strings"

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
		res, err := chat.CreateCompletion(gateway.CreateCompletionRequest{Description: composeCLICompletionPrompt(desc)})
		if err != nil {
			fmt.Println(err)
			return
		}

		suggestions := ParseCompletion(res.Choices[0].Messages.Content)

		for _, s := range suggestions {
			fmt.Println("Command: ", s.Command)
			fmt.Println("Note: ", s.Note)
			fmt.Println()
		}
	},
}

type Suggestion struct {
	Command string
	Note    string
}

func ParseCompletion(res string) []Suggestion {
	suggestions := []Suggestion{}

	cur := Suggestion{}
	for _, line := range strings.Split(res, "\n") {
		if strings.HasPrefix(line, "command:") {
			cur.Command = strings.TrimSpace(strings.TrimPrefix(line, "command:"))
		}
		if strings.HasPrefix(line, "note:") {
			cur.Note = strings.TrimSpace(strings.TrimPrefix(line, "note:"))
		}
		if cur.Command != "" && cur.Note != "" {
			suggestions = append(suggestions, cur)
			cur = Suggestion{}
		}
	}

	return suggestions
}

func composeCLICompletionPrompt(desc string) string {
	return `What CLI command will accomplish the following objectives? The answer should follow this format 'command:XXX
note:XXX'. Objectives: ` + desc
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

}
