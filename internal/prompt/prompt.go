package prompt

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/manifoldco/promptui"
)

type SelectItem struct {
	Command     string
	Summary     string
	Description string
}

func ToSelectItems(res string) []SelectItem {
	items := []SelectItem{}

	cur := SelectItem{}
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
			items = append(items, cur)
			cur = SelectItem{}
		}
	}

	return items
}

func SelectCommand(stdin io.ReadCloser, items []SelectItem) (int, error) {
	if len(items) == 0 {
		return 0, fmt.Errorf("no items to select")
	}

	prompt := promptui.Select{
		Stdin: stdin,
		Label: "Select command",
		Items: items,
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
		}}

	i, _, err := prompt.Run()

	if err != nil {
		return i, fmt.Errorf("failed to run prompt: %w", err)
	}
	return i, nil
}

func ParseToChatGPTInput(desc string) string {
	return `Please provide up to 3 commands that accomplish the following objectives. Each candidate must follow the format "Command: XXX
Summary: XXX
Description: XXX" and output them consecutively to form a single answer. Objectives: ` + desc
}

func PromptApiKey() (string, error) {
	prompt := promptui.Prompt{
		Label: "Please enter your OpenAI API key",
		Validate: func(input string) error {
			if len(input) == 0 {
				return errors.New("API key cannot be empty")
			}
			return nil
		},
		Mask: '*',
	}

	return prompt.Run()
}

type CommandHandlers struct {
	OnRun  func() error
	OnCopy func() error
	OnExit func() error
}

func SelectCommandActions(stdin io.ReadCloser, handlers CommandHandlers) error {
	items := []struct {
		Label   string
		Handler func() error
	}{
		{Label: "Run command", Handler: handlers.OnRun},
		{Label: "Copy command", Handler: handlers.OnCopy},
		{Label: "Exit", Handler: handlers.OnExit},
	}

	prompt := promptui.Select{
		Stdin: stdin,
		Label: "What do you want to do?",
		Items: items,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ .Label }}?",
			Active:   "👉 {{ .Label | cyan }}",
			Inactive: "   {{ .Label | cyan }}",
			Selected: "👉 {{ .Label | cyan }}",
		}}

	i, _, err := prompt.Run()

	if err != nil {
		return fmt.Errorf("failed to run prompt: %w", err)
	}

	return items[i].Handler()
}
