package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/kj455/intelli-cli/gateway"
	"github.com/kj455/intelli-cli/utils"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "?",
	Short: "IntelliCLI helps you find exact command you want to run",
	Long:  `IntelliCLI to help you find exact command you want to run`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		desc := args[0]

		res, err := utils.WithLoading(os.Stdout, "🤔 Thinking...", func() (gateway.CreateCompletionResponse, error) {
			chat := gateway.CreateChatClient()
			res, err := chat.CreateCompletion(gateway.CreateCompletionRequest{Description: composeCLICompletionPrompt(desc)})

			if err != nil {
				return res, fmt.Errorf("failed to create completion: %w", err)
			}
			return res, nil
		})

		if err != nil {
			fmt.Println(err)
			return
		}

		suggestions := ParseCompletion(res.Choices[0].Messages.Content)

		prompt := promptui.Select{
			Label: "👇 Select a command",
			Items: suggestions,
			Templates: &promptui.SelectTemplates{
				Label:    "{{ . }}?",
				Active:   "✔︎ {{ .Command | cyan }}",
				Inactive: "  {{ .Command | cyan }}",
				Selected: "✔︎ {{ .Command | cyan }}",
				Details: `
--------- Command ----------
{{ "Command:" | faint }}	{{ .Command }}
{{ "Description:" | faint }}	{{ .Description }}
`,
			},
		}

		i, _, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		result, err := exec.Command("bash", "-c", suggestions[i].Command).Output()
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(string(result))
	},
}

type Suggestion struct {
	Command     string
	Summary     string
	Description string
}

func ParseCompletion(res string) []Suggestion {
	suggestions := []Suggestion{}

	cur := Suggestion{}
	for _, line := range strings.Split(res, "\n") {
		if strings.HasPrefix(line, "Command:") {
			cur.Command = strings.TrimSpace(strings.TrimPrefix(line, "Command:"))
		}
		if strings.HasPrefix(line, "Summary:") {
			cur.Summary = strings.TrimSpace(strings.TrimPrefix(line, "Summary:"))
		}
		if strings.HasPrefix(line, "Description:") {
			cur.Description = strings.TrimSpace(strings.TrimPrefix(line, "Description:"))
		}
		if cur.Command != "" && cur.Description != "" {
			suggestions = append(suggestions, cur)
			cur = Suggestion{}
		}
	}

	return suggestions
}

func composeCLICompletionPrompt(desc string) string {
	return `Please provide up to 3 CLI commands that accomplish the following objectives Each candidate should follow the format "Command: XXX
Summary: XXX
Description: XXX" and output them consecutively to form a single answer. Objectives: ` + desc
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
