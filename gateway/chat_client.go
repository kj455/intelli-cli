package gateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/kj455/intelli-cli/secret"
)

type CreateCompletionRequest struct {
	Description string
}
type CreateCompletionResponse ApiCreateChatCompletionResponse

type ChatClient struct {
	client HTTPClient
}

type ChatClientInterface interface {
	CreateCompletion(payload CreateCompletionRequest) (CreateCompletionResponse, error)
}

func (c *ChatClient) CreateCompletion(payload CreateCompletionRequest) (CreateCompletionResponse, error) {
	var res CreateCompletionResponse

	jsonBody, err := json.Marshal(&ApiCreateChatCompletionRequest{
		Model: "gpt-3.5-turbo",
		Messages: []Message{
			{
				Role:    "user",
				Content: payload.Description,
			},
		},
	})
	if err != nil {
		return res, fmt.Errorf("failed to marshal json: %w", err)
	}

	req, err := http.NewRequest("POST", BASE_URL+"/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		return res, fmt.Errorf("failed to create request: %w", err)
	}
	AuthorizeRequest(req)

	resp, err := c.client.Do(req)

	if resp.StatusCode != http.StatusOK || err != nil {
		if resp.StatusCode == http.StatusUnauthorized {
			secret.DeleteApiKey()
			return res, fmt.Errorf("unauthorized, please check your api key")
		}
		return res, fmt.Errorf("failed to do request: %w", err)
	}
	defer resp.Body.Close()

	resBody, _ := ioutil.ReadAll(resp.Body)
	if err = json.Unmarshal(resBody, &res); err != nil {
		return res, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return res, nil
}

func NewChatClient(httpClient HTTPClient) *ChatClient {
	return &ChatClient{
		client: httpClient,
	}
}

func CreateChatClient() *ChatClient {
	return NewChatClient(&http.Client{})
}
