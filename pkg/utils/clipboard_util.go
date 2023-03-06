package utils

import "golang.design/x/clipboard"

func CopyToClipboard(content string) error {
	err := clipboard.Init()
	if err != nil {
		return err
	}

	clipboard.Write(clipboard.FmtText, []byte(content))
	return nil
}
