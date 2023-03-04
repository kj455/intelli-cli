package gateway

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type CreateCompletionRequest struct {
	Description string
}
type CreateCompletionResponse ApiCreateChatCompletionResponse

type ChatClient struct {
	client HTTPClient
}

func (c *ChatClient) CreateCompletion(payload CreateCompletionRequest) (ApiCreateChatCompletionResponse, error) {
	jsonBody, err := json.Marshal(&ApiCreateChatCompletionRequest{
		Model: "gpt-3.5-turbo",
		Messages: []Message{
			{
				Role:    "user",
				Content: "Please return only the CLI command string for the following purposes." + payload.Description,
			},
		},
	})
	if err != nil {
		return ApiCreateChatCompletionResponse{}, err
	}

	req, err := CreateHttpRequest("POST", "/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		return ApiCreateChatCompletionResponse{}, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return ApiCreateChatCompletionResponse{}, err
	}
	defer resp.Body.Close()

	resBody, _ := ioutil.ReadAll(resp.Body)

	var res ApiCreateChatCompletionResponse
	if err = json.Unmarshal(resBody, &res); err != nil {
		return ApiCreateChatCompletionResponse{}, err
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
