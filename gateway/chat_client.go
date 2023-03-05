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
		return res, err
	}

	req, err := CreateHttpRequest("POST", "/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		return res, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return res, err
	}
	defer resp.Body.Close()

	resBody, _ := ioutil.ReadAll(resp.Body)

	if err = json.Unmarshal(resBody, &res); err != nil {
		return res, err
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
