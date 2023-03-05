package gateway

import (
	"fmt"
	"io"
	"net/http"

	"github.com/kj455/intelli-cli/secret"
)

const BASE_URL = "https://api.openai.com/v1"

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func CreateHttpRequest(method string, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, BASE_URL+path, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	apiKey, err := secret.GetApiKey()
	if err != nil {
		return nil, fmt.Errorf("failed to get api key: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	return req, nil
}
