package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type LLMClient struct {
	provider   string
	apiKey     string
	httpClient *http.Client
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionRequest struct {
	Model     string    `json:"model"`
	MaxTokens int       `json:"max_tokens"`
	Messages  []Message `json:"messages"`
}

type MessageResponseContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type UsageResponse struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

type ChatCompletionResponse struct {
	Content      []MessageResponseContent `json:"content"`
	Id           string                   `json:"id"`
	Model        string                   `json:"model"`
	Role         string                   `json:"role"`
	StopReason   string                   `json:"stop_reason"`
	StopSequence int                      `json:"stop_sequence"`
	Type         string                   `json:"type"`
	Usage        UsageResponse            `json:"usage"`
}

func NewLLMClient(provider string, apiKey string) *LLMClient {
	return &LLMClient{
		provider:   provider,
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}
}

func (c *LLMClient) SendMessage(content string) (*ChatCompletionResponse, error) {
	request := ChatCompletionRequest{
		Model:     "claude-3-5-sonnet-20241022",
		MaxTokens: 1024,
		Messages: []Message{
			{
				Role:    "user",
				Content: content,
			},
		},
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error building JSON: %w", err)
	}

	req, err := http.NewRequest(
		"POST",
		"https://api.anthropic.com/v1/messages",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("error building request: %w", err)

	}

	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("content-type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error during request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error (status: %d): %s", resp.StatusCode, string(body))
	}

	var apiResp ChatCompletionResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	return &apiResp, nil
}
