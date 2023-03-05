package utils

import (
	"fmt"
	"io"
)

func WithLoading[T any](w io.Writer, loading string, f func() (T, error)) (T, error) {
	fmt.Fprint(w, loading)
	res, err := f()
	fmt.Fprint(w, "\033[2K\r")
	return res, err
}
