package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/kj455/intelli-cli/gateway"
	"github.com/kj455/intelli-cli/utils"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var Dry bool

var rootCmd = &cobra.Command{
	Use:   "intelli-cli",
	Short: "IntelliCLI helps you find exact command you want to run",
	Long:  `IntelliCLI to help you find exact command you want to run`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunRoot(Context{
			stdin:  os.Stdin,
			stdout: os.Stdout,
			chat:   gateway.CreateChatClient(),
		}, args)
	},
}

type Context struct {
	stdin  io.ReadCloser
	stdout io.Writer
	chat   gateway.ChatClientInterface
}

func RunRoot(ctx Context, args []string) error {
	desc := args[0]

	res, err := utils.WithLoading(ctx.stdout, "🤔 Thinking...", func() (gateway.CreateCompletionResponse, error) {
		res, err := ctx.chat.CreateCompletion(gateway.CreateCompletionRequest{Description: composeCLICompletionPrompt(desc)})

		if err != nil {
			return res, fmt.Errorf("failed to create completion: %w", err)
		}
		return res, nil
	})

	if err != nil {
		return err
	}

	if len(res.Choices) == 0 {
		return fmt.Errorf("no suggestions found")
	}
	suggestions := ToSuggestions(res.Choices[0].Messages.Content)

	if len(suggestions) == 0 {
		return fmt.Errorf("no suggestions found")
	}

	prompt := promptui.Select{
		Stdin: ctx.stdin,
		Label: "Select a command",
		Items: suggestions,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}?",
			Active:   "👉 {{ .Command | cyan }}",
			Inactive: "   {{ .Command | cyan }}",
			Selected: "👉 {{ .Command | cyan }}",
			Details: `
--------- Command ----------
{{ "Command:" | faint }}	{{ .Command }}
{{ "Description:" | faint }}	{{ .Description }}
`,
		},
	}

	i, _, err := prompt.Run()
	if err != nil {
		return fmt.Errorf("failed to run prompt: %w", err)
	}

	if Dry {
		fmt.Fprintln(ctx.stdout, "😌 Dry run, skipping command")
		return nil
	}

	command := suggestions[i].Command

	fmt.Fprintln(ctx.stdout, "🚀 Run command: "+command)

	result, err := exec.Command("bash", "-c", command).Output()
	if err != nil {
		return fmt.Errorf("failed to run command: %w", err)
	}

	fmt.Fprintln(ctx.stdout, string(result))

	return nil
}

type Suggestion struct {
	Command     string
	Summary     string
	Description string
}

func ToSuggestions(res string) []Suggestion {
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
	return `Please provide up to 3 commands that accomplish the following objectives Each candidate should strictly follow the format "Command: XXX
Summary: XXX
Description: XXX" and output them consecutively to form a single answer. Objectives: ` + desc
}

func init() {
	rootCmd.Flags().BoolVarP(&Dry, "dry", "d", false, "dry run")
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
