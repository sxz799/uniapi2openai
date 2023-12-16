package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	myModel "github/sxz799/gemini2chatgpt/model"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"log"
	"time"
)

func DoTrans(apiKey string, openaiBody myModel.ChatGPTRequestBody, c *gin.Context) {
	ctx := context.Background()
	// Access your API key as an environment variable (see "Set up your API key" above)
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// For text-only input, use the gemini-pro model
	model := client.GenerativeModel("gemini-pro")
	// Initialize the chat
	cs := model.StartChat()
	cs.History = []*genai.Content{}
	lastMsg := ""
	for i, msg := range openaiBody.Messages {
		content := msg.Content
		role := msg.Role

		if i == len(openaiBody.Messages)-1 {
			lastMsg = content
			break
		}
		if role == "assistant" {
			role = "model"
		}
		cs.History = append(cs.History,
			&genai.Content{
				Parts: []genai.Part{
					genai.Text(content),
				},
				Role: role,
			},
		)
	}

	iter := cs.SendMessageStream(ctx, genai.Text(lastMsg))
	for {
		resp, err := iter.Next()
		if err == iterator.Done {
			c.Writer.WriteString("data: [DONE]\n")
			c.Writer.Flush()
			break
		}
		if err != nil {
			break
		}
		for _, candidate := range resp.Candidates {
			for _, p := range candidate.Content.Parts {
				str := fmt.Sprintf("%s", p)
				dd := myModel.Delta{
					Role:    "assistant",
					Content: str,
				}
				cc := myModel.ChoiceChunk{
					Delta: dd,
				}
				chunk := myModel.ChatCompletionChunk{
					ID:      fmt.Sprintf("%d", time.Now().Unix()),
					Object:  "chat.completion.chunk",
					Created: time.Now().Unix(),
					Model:   "gemini-pro",
					Choices: []myModel.ChoiceChunk{cc},
				}
				marshal, _ := json.Marshal(chunk)
				_, err = c.Writer.WriteString("data: " + string(marshal) + "\n")
				if err != nil {
					return
				}
				c.Writer.Flush()
			}
		}

	}
}
