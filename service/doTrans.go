package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"github/sxz799/gemini2chatgpt/model"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"log"
	"strings"
)

func DoTrans(apiKey string, openaiBody model.ChatGPTRequestBody, c *gin.Context) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	gModel := client.GenerativeModel("gemini-pro")
	cs := gModel.StartChat()
	cs.History = []*genai.Content{}
	lastMsg := ""
	for i, msg := range openaiBody.Messages {
		content := msg.Content
		role := strings.ReplaceAll(msg.Role, "assistant", "model")
		if i == len(openaiBody.Messages)-1 {
			lastMsg = content
			break
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

	if openaiBody.Stream {
		SendStreamResponse(cs, ctx, lastMsg, c)
	} else {
		SendSingleResponse(cs, ctx, lastMsg, c)
	}
}

func SendStreamResponse(cs *genai.ChatSession, ctx context.Context, lastMsg string, c *gin.Context) {
	iter := cs.SendMessageStream(ctx, genai.Text(lastMsg))
	for {
		resp, err := iter.Next()
		if errors.Is(err, iterator.Done) {
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
				chunk := model.NewChatCompletionChunk(str, "gemini-pro")
				marshal, _ := json.Marshal(chunk)
				_, err = c.Writer.WriteString("data: " + string(marshal) + "\n\n")
				if err != nil {
					return
				}
				c.Writer.Flush()
			}
		}
	}
}

func SendSingleResponse(cs *genai.ChatSession, ctx context.Context, lastMsg string, c *gin.Context) {
	resp, err := cs.SendMessage(ctx, genai.Text(lastMsg))
	if err!=nil{
		c.String(200, "SendMessage Error:",err.Error())
		return
	}
	if len(resp.Candidates) < 1 || len(resp.Candidates[0].Content.Parts) < 1 {
		c.String(200, "no response")
		return
	}
	part := resp.Candidates[0].Content.Parts[0]
	str := fmt.Sprintf("%s", part)
	cc := model.NewChatCompletion(str, "gemini-pro")
	marshal, _ := json.Marshal(cc)
	_, err = c.Writer.Write(marshal)
	if err != nil {
		return
	}
	c.Writer.Flush()
}
