package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github/sxz799/gemini2chatgpt/model"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

var IngoreSystemPrompt bool

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
	if len(openaiBody.Messages) < 1 {
		c.JSON(200, gin.H{
			"err": "会话不可为空",
		})
	}
	var lastMsg string
	var lastRole string
	for i, msg := range openaiBody.Messages {
		if msg.Role == "system" {
			if IngoreSystemPrompt {
				openaiBody.Messages[i].Content = "你好!"
			}
			openaiBody.Messages[i].Role = "user"
		}
		// 将assistant角色替换为model
		if msg.Role == "assistant" {
			openaiBody.Messages[i].Role = "model"
		}
	}
	for i, msg := range openaiBody.Messages {
		if i == 0 && msg.Role != "user" {
			c.JSON(200, gin.H{
				"err": "第一条会话必须是用户发起",
			})
			return
		}
		if msg.Role == "user" {
			lastMsg = lastMsg + msg.Content
		} else {
			lastMsg = ""
		}
		if i == len(openaiBody.Messages)-1 {
			break
		}
		if msg.Role != lastRole {
			cs.History = append(cs.History, &genai.Content{Parts: []genai.Part{genai.Text(msg.Content)}, Role: msg.Role})
		} else {
			cs.History[len(cs.History)-1].Parts = append(cs.History[len(cs.History)-1].Parts, []genai.Part{genai.Text(msg.Content)}...)
		}
		lastRole = msg.Role
	}

	if len(cs.History) > 1 && cs.History[len(cs.History)-1].Role == "user" {
		cs.History = cs.History[:len(cs.History)-1]
	}
	if len(cs.History) == 1 {
		cs.History = []*genai.Content{}
	}

	fmt.Println("======历史记录======")
	for _, hs := range cs.History {
		fmt.Println(hs.Role, ":", hs.Parts)
	}
	fmt.Println("====================")
	fmt.Println("user:", lastMsg)
	fmt.Println("====================")

	if openaiBody.Stream {
		//支持 SSE特性
		c.Writer.Header().Set("Transfer-Encoding", "chunked")
		c.Writer.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
		sendStreamResponse(cs, ctx, lastMsg, c)
	} else {
		sendSingleResponse(cs, ctx, lastMsg, c)
	}
}

func sendStreamResponse(cs *genai.ChatSession, ctx context.Context, lastMsg string, c *gin.Context) {
	iter := cs.SendMessageStream(ctx, genai.Text(lastMsg))
	for {
		id := fmt.Sprintf("chatcmpl-%d", time.Now().Unix())
		resp, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			str := "stop"
			cc := model.ChoiceChunk{
				FinishReason: &str,
				Index:        0,
			}
			tChatCompletionChunk := model.ChatCompletionChunk{
				ID:      id,
				Model:   "gemini-pro",
				Created: time.Now().Unix(),
				Object:  "chat.completion.chunk",
				Choices: []model.ChoiceChunk{cc},
			}

			marshal, _ := json.Marshal(tChatCompletionChunk)
			_, _ = c.Writer.WriteString(fmt.Sprintf("data: %s\n\n", marshal))
			_, _ = c.Writer.WriteString("data: [DONE]\n")
			c.Writer.Flush()
			break
		}
		if err != nil {
			c.JSON(200, gin.H{
				"lastMsg": lastMsg,
				"err":     err.Error(),
			})
			break
		}

		for _, candidate := range resp.Candidates {
			for _, p := range candidate.Content.Parts {
				str := fmt.Sprintf("%s", p)
				chunk := model.NewChatCompletionChunk(id, str, "gemini-pro")
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

func sendSingleResponse(cs *genai.ChatSession, ctx context.Context, lastMsg string, c *gin.Context) {
	resp, err := cs.SendMessage(ctx, genai.Text(lastMsg))
	if err != nil {
		c.String(200, "SendMessage Error:", err.Error())
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
