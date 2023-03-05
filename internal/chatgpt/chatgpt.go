package chatgpt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const BASE_URL = "https://api.openai.com/v1"

type ApiCreateChatCompletionRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type ApiCreateChatCompletionResponse struct {
	Id      string   `json:"id"`
	Object  string   `json:"object"`
	Created int      `json:"created"`
	Model   string   `json:"model"`
	Usage   Usage    `json:"usage"`
	Choices []Choice `json:"choices"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type Choice struct {
	Index        int     `json:"index"`
	Messages     Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type ChatGPT struct {
	apiKey string
}

func New(apiKey string) *ChatGPT {
	return &ChatGPT{apiKey: apiKey}
}

type CreateChatCompletionRequest struct {
	Description string
}
type CreateChatCompletionResponse ApiCreateChatCompletionResponse

func (c *ChatGPT) CreateCompletion(payload CreateChatCompletionRequest) (CreateChatCompletionResponse, error) {
	var res CreateChatCompletionResponse

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
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := http.DefaultClient.Do(req)

	if resp.StatusCode != http.StatusOK || err != nil {
		if resp.StatusCode == http.StatusUnauthorized {
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
