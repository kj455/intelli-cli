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
	return `Please provide up to 3 commands that accomplish the following objectives. Each candidate should strictly follow the format "Command: XXX
Summary: XXX
Description: XXX" and output them consecutively to form a single answer. Objectives: ` + desc
}

func PropmtApiKey() (string, error) {
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
