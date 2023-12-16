package model

import (
	"fmt"
	"time"
)

type ChatGPTRequestBody struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Delta struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Choice struct {
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
	Index        int     `json:"index"`
}
type ChoiceChunk struct {
	Delta        Delta  `json:"delta"`
	FinishReason string `json:"finish_reason"`
	Index        int    `json:"index"`
}

type ChatCompletion struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Usage   Usage    `json:"usage"`
	Choices []Choice `json:"choices"`
}

type ChatCompletionChunk struct {
	ID      string        `json:"id"`
	Object  string        `json:"object"`
	Created int64         `json:"created"`
	Model   string        `json:"model"`
	Choices []ChoiceChunk `json:"choices"`
}

func NewChatCompletionChunk(msg, model string) *ChatCompletionChunk {
	cc := ChoiceChunk{
		Delta: Delta{
			Role:    "assistant",
			Content: msg,
		},
	}
	return &ChatCompletionChunk{
		ID:      fmt.Sprintf("%d", time.Now().Unix()),
		Object:  "chat.completion.chunk",
		Created: time.Now().Unix(),
		Model:   model,
		Choices: []ChoiceChunk{cc},
	}
}

func NewChatCompletion(msg, model string) *ChatCompletion {
	cho := Choice{
		Message: Message{
			Role:    "assistant",
			Content: msg,
		},
	}
	return &ChatCompletion{
		ID:      fmt.Sprintf("%d", time.Now().Unix()),
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   model,
		Choices: []Choice{cho},
	}
}
