package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/kj455/intelli-cli/internal/chatgpt"
	"github.com/kj455/intelli-cli/internal/prompt"
	"github.com/kj455/intelli-cli/internal/secret"
	"github.com/kj455/intelli-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var questionCmd = &cobra.Command{
	Use:   "q",
	Short: "q helps you find exact command you want to run",
	Long:  `q helps you find exact command you want to run`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := secret.SetupSecretIfNeeded()
		if err != nil {
			return fmt.Errorf("failed to setup secret. Please exec `intelli-cli auth refresh`")
		}

		key, err := secret.GetApiKey()
		if err != nil {
			return fmt.Errorf("failed to get api key. Please exec `intelli-cli auth refresh`")
		}

		return RunQuestionCmd(Context{
			stdin:  os.Stdin,
			stdout: os.Stdout,
			stderr: os.Stderr,
			chat:   *chatgpt.New(key),
		}, args)
	},
}

type Context struct {
	stdin  io.ReadCloser
	stdout io.Writer
	stderr io.Writer
	chat   chatgpt.ChatGPT
}

func RunQuestionCmd(ctx Context, args []string) error {
	desc := prompt.ParseToChatGPTInput(args[0])

	res, err := utils.WithLoading(ctx.stdout, "🤔 Thinking...", func() (chatgpt.CreateChatCompletionResponse, error) {
		res, err := ctx.chat.CreateCompletion(chatgpt.CreateChatCompletionRequest{Description: desc})

		if err != nil {
			return res, fmt.Errorf("failed to create completion: %w", err)
		}
		return res, nil
	})

	if err != nil {
		return err
	}

	if len(res.Choices) == 0 {
		return fmt.Errorf("no items found")
	}

	items := prompt.ToSelectItems(res.Choices[0].Messages.Content)
	item, err := prompt.SelectCommand(ctx.stdin, items)
	command := item.Command

	if err != nil {
		return fmt.Errorf("failed to run prompt: %w", err)
	}

	prompt.SelectCommandActions(ctx.stdin, prompt.CommandHandlers{
		OnRun: func() error {
			fmt.Fprintln(ctx.stdout, "🚀 Running: "+command)

			res, err := utils.RunCommand(command)
			if err != nil {
				return fmt.Errorf("failed to run command: %w", err)
			}

			fmt.Fprintln(ctx.stdout, res)
			return nil
		},
		OnExit: func() error {
			fmt.Fprintln(ctx.stdout, "👋 Bye")
			return nil
		},
	})

	return nil
}

func init() {}
