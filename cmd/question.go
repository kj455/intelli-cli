package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/kj455/intelli-cli/internal/chatgpt"
	"github.com/kj455/intelli-cli/internal/prompt"
	"github.com/kj455/intelli-cli/internal/secret"
	"github.com/kj455/intelli-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var Dry bool

var questionCmd = &cobra.Command{
	Use:   "q",
	Short: "q helps you find exact command you want to run",
	Long:  `q to help you find exact command you want to run`,
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

		return RunRoot(Context{
			stdin:  os.Stdin,
			stdout: os.Stdout,
			chat:   *chatgpt.New(key),
		}, args)
	},
}

type Context struct {
	stdin  io.ReadCloser
	stdout io.Writer
	chat   chatgpt.ChatGPT
}

func RunRoot(ctx Context, args []string) error {
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
	i, err := prompt.SelectCommand(ctx.stdin, items)

	if err != nil {
		return fmt.Errorf("failed to run prompt: %w", err)
	}

	if Dry {
		fmt.Fprintln(ctx.stdout, "😌 Dry run, skipping command")
		return nil
	}

	command := items[i].Command

	fmt.Fprintln(ctx.stdout, "🚀 Run command: "+command)

	result, err := exec.Command("bash", "-c", command).Output()
	if err != nil {
		return fmt.Errorf("failed to run command: %w", err)
	}

	fmt.Fprintln(ctx.stdout, string(result))

	return nil
}

func init() {
	questionCmd.Flags().BoolVarP(&Dry, "dry", "d", false, "dry run")
}
