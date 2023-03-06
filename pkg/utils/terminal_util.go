package utils

import (
	"fmt"
	"io"
	"os/exec"
)

const (
	CLEAR_LINE = "\033[2K\r"
)

func WithLoading[T any](w io.Writer, loading string, f func() (T, error)) (T, error) {
	fmt.Fprint(w, loading)
	res, err := f()
	fmt.Fprint(w, CLEAR_LINE)
	return res, err
}

func RunCommand(command string) (string, error) {
	cmd := exec.Command("bash", "-c", command)

	res, err := cmd.CombinedOutput()

	if err != nil {
		return "", fmt.Errorf("output: %s, error: %w", string(res), err)
	}

	return string(res), nil
}
