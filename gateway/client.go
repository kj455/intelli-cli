package gateway

import (
	"io"
	"net/http"
	"os"
)

const BASE_URL = "https://api.openai.com/v1"

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func CreateHttpRequest(method string, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, BASE_URL+path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))

	return req, nil
}
