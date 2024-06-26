package u2o4gemini

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sxz799/uniapi2openai/config"
	"github.com/sxz799/uniapi2openai/model"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func DoTrans(ignoreSystemPrompt bool, openaiBody model.OpenaiBody, c *gin.Context) {
	key := c.GetHeader("Authorization")
	if len(strings.Split(key, " ")) != 2 {
		if key == "" {
			c.JSON(400, gin.H{
				"error": "Authorization header is invalid",
			})
			return
		}
	} else {
		key = strings.Split(key, " ")[1]
	}
	ctx := context.Background()
	clientOptionApi := option.WithAPIKey(key)
	customUrl := config.GeminiProxyUrl
	if customUrl == "" {
		customUrl = "https://generativelanguage.googleapis.com"
	}
	clientOptionUrl := option.WithEndpoint(customUrl)
	client, err := genai.NewClient(ctx, clientOptionApi, clientOptionUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	modelName := openaiBody.Model
	gModel := client.GenerativeModel(modelName)
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
			if ignoreSystemPrompt {
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

	if openaiBody.Stream {
		//支持 SSE特性
		c.Writer.Header().Set("Transfer-Encoding", "chunked")
		c.Writer.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
		sendStreamResponse(cs, ctx, lastMsg, modelName, c)
	} else {
		sendSingleResponse(cs, ctx, lastMsg, modelName, c)
	}
}

func sendStreamResponse(cs *genai.ChatSession, ctx context.Context, lastMsg, modelName string, c *gin.Context) {
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
				Model:   modelName,
				Created: time.Now().Unix(),
				Object:  "chat.completion.chunk",
				Choices: []model.ChoiceChunk{cc},
			}

			marshal, _ := json.Marshal(tChatCompletionChunk)
			_, _ = c.Writer.WriteString(fmt.Sprintf("data: %s\n\n", marshal))
			c.Writer.Flush()
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
				chunk := model.NewChatCompletionChunk(id, str, modelName)
				marshal, _ := json.Marshal(chunk)
				_, err = c.Writer.WriteString("data: " + string(marshal) + "\n\n")
				if err != nil {
					break
				}
				c.Writer.Flush()
			}
		}
	}
}

func sendSingleResponse(cs *genai.ChatSession, ctx context.Context, lastMsg, modelName string, c *gin.Context) {
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
	cc := model.NewChatCompletion(str, modelName)
	marshal, _ := json.Marshal(cc)
	_, err = c.Writer.Write(marshal)
	if err != nil {
		return
	}
	c.Writer.Flush()
}
