package u2o4tongyi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sxz799/uniapi2openai/model"
	"net/http"
	"strings"
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
	tongYiBody := transOpenAIReq2TongYiReq(ignoreSystemPrompt, openaiBody)
	jsonData, _ := json.Marshal(tongYiBody)
	tyApi := "https://dashscope.aliyuncs.com/api/v1/services/aigc/text-generation/generation"

	req, _ := http.NewRequest("POST", tyApi, bytes.NewReader(jsonData))

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-DashScope-SSE", "enable")
	req.Header.Add("Authorization", "Bearer "+key)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	c.Writer.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	msgChan := make(chan string, 10)
	go func() {
		defer close(msgChan)
		id := ""
		for {
			buf := make([]byte, 4096)
			n, err2 := resp.Body.Read(buf)
			if err2 != nil {
				tChunk := model.NewStopChatCompletionChunk(id, openaiBody.Model)
				tMarshal, _ := json.Marshal(tChunk)
				msgChan <- fmt.Sprintf("data: %s\n\n", tMarshal)
				msgChan <- fmt.Sprintf("data: [DONE]\n")
				break
			}
			var str = string(buf[:n])
			//截取data:之后的内容
			index := strings.Index(str, "{")
			if index > 0 {
				str = str[index:]
			}
			result, tid, finish := transTongYiResp2OpenAIResp(openaiBody.Model, str)
			if id == "" {
				id = tid
			}
			msgChan <- fmt.Sprintf("data: %s\n\n", result)
			if finish {
				chunk := model.NewStopChatCompletionChunk(id, openaiBody.Model)
				marshal, _ := json.Marshal(chunk)
				msgChan <- fmt.Sprintf("data: %s\n\n", marshal)
				msgChan <- fmt.Sprintf("data: [DONE]\n")
				break
			}
		}
	}()
	for msg := range msgChan {
		_, _ = c.Writer.WriteString(msg)
		c.Writer.Flush()
	}
}

func transOpenAIReq2TongYiReq(ignoreSystemPrompt bool, body model.OpenaiBody) *TongyiBody {
	modelName := body.Model
	tMessages := body.Messages
	var messages []Message
	lastRole := ""
	for _, message := range tMessages {
		if ignoreSystemPrompt && message.Role == "system" {
			continue
		}
		if message.Role == lastRole {
			messages[len(messages)-1].Content += message.Content
			continue
		}
		messages = append(messages, Message{
			Role:    message.Role,
			Content: message.Content,
		})
		lastRole = message.Role
	}
	return NewTongyiBody(modelName, messages)
}

func transTongYiResp2OpenAIResp(modelName, origin string) (result, id string, finish bool) {

	var qwResp TongYiResponseBody
	err := json.Unmarshal([]byte(origin), &qwResp)
	if err != nil {
		return
	}
	if len(qwResp.Output.Choices) < 1 {
		return
	}
	//fmt.Println("====================")
	//fmt.Println(qwResp.Output.Choices[0].Message.Content)
	//fmt.Println("====================")
	chunk := model.NewChatCompletionChunk(qwResp.RequestID, qwResp.Output.Choices[0].Message.Content, modelName)
	marshal, _ := json.Marshal(chunk)
	result = string(marshal)
	id = qwResp.RequestID
	if qwResp.Output.Choices[0].FinishReason == "stop" {
		finish = true
	}
	return
}
