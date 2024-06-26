package model

import (
	"fmt"
	"time"
)

type OpenaiBody struct {
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
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

type Choice struct {
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
	Index        int     `json:"index"`
}
type ChoiceChunk struct {
	Delta        Delta   `json:"delta"`
	FinishReason *string `json:"finish_reason"`
	Index        int     `json:"index"`
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

// NewChatCompletionChunk 流式输出 chunk
func NewChatCompletionChunk(id, msg, model string) *ChatCompletionChunk {
	cc := ChoiceChunk{
		Delta: Delta{
			Role:    "assistant",
			Content: msg,
		},
		Index: 0,
	}
	return &ChatCompletionChunk{
		ID:      id,
		Object:  "chat.completion.chunk",
		Created: time.Now().Unix(),
		Model:   model,
		Choices: []ChoiceChunk{cc},
	}
}

// NewChatCompletion 单次输出
func NewChatCompletion(msg, model string) *ChatCompletion {
	cho := Choice{
		Message: Message{
			Role:    "assistant",
			Content: msg,
		},
	}
	return &ChatCompletion{
		ID:      fmt.Sprintf("chatcmpl-%d", time.Now().Unix()),
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   model,
		Choices: []Choice{cho},
	}
}

func NewStopChatCompletionChunk(id, modelName string) *ChatCompletionChunk {
	tStr := "stop"
	cc := ChoiceChunk{
		FinishReason: &tStr,
		Index:        0,
	}
	tChatCompletionChunk := ChatCompletionChunk{
		ID:      id,
		Model:   modelName,
		Created: time.Now().Unix(),
		Object:  "chat.completion.chunk",
		Choices: []ChoiceChunk{cc},
	}
	return &tChatCompletionChunk
}
