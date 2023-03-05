package utils

import (
	"fmt"
	"io"
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
