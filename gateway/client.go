package gateway

import (
	"fmt"
	"net/http"

	"github.com/kj455/intelli-cli/secret"
)

const BASE_URL = "https://api.openai.com/v1"

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func AuthorizeRequest(req *http.Request) error {
	apiKey, err := secret.GetApiKey()
	if err != nil {
		return fmt.Errorf("failed to get api key: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	return nil
}
