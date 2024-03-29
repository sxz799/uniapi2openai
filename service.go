package main

import (
	"fmt"
	"time"
)


func newChatCompletionChunk(id, msg, model string) *ChatCompletionChunk {
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

func newChatCompletion(msg, model string) *ChatCompletion {
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