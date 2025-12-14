package ai

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

type Client struct {
	client *anthropic.Client
	model  anthropic.Model
}

// NewClient creates a new Anthropic API client
func NewClient(apiKey string, model string) *Client {
	client := anthropic.NewClient(
		option.WithAPIKey(apiKey),
	)

	return &Client{
		client: &client,
		model:  anthropic.Model(model),
	}
}

// SendMessage sends a message to Claude and returns the response
func (c *Client) SendMessage(ctx context.Context, system string, user string) (string, error) {
	message, err := c.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     c.model,
		MaxTokens: 2048,
		System: []anthropic.TextBlockParam{
			{
				Text: system,
			},
		},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(user)),
		},
	})
	if err != nil {
		return "", err
	}

	// Extract text from the response by marshaling and unmarshaling
	if len(message.Content) > 0 {
		// Use JSON to extract the text field
		data, err := json.Marshal(message.Content[0])
		if err == nil {
			var result struct {
				Text string `json:"text"`
			}
			if err := json.Unmarshal(data, &result); err == nil && result.Text != "" {
				return result.Text, nil
			}
		}

		return "", fmt.Errorf("unable to extract text from response")
	}

	return "", fmt.Errorf("no content in response")
}
