package u2o4tongyiWeb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sxz799/uniapi2openai/model"
	"net/http"
	"strings"
	"time"
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

	//1. 准备发起请求的请求头
	var headers = fakeHeader
	cookieStr := generateCookie(key)
	headers["Cookie"] = cookieStr
	//2. 准备发起请求的请求体

	messages := openaiBody.Messages

	reqBodyMap := make(map[string]interface{})
	reqBodyMap["model"] = ""
	reqBodyMap["action"] = "next"
	reqBodyMap["mode"] = "chat"
	reqBodyMap["userAction"] = "chat"
	reqBodyMap["requestId"] = uuid.NewString()
	reqBodyMap["sessionId"] = ""
	reqBodyMap["sessionType"] = "text_chat"
	reqBodyMap["parentMsgId"] = ""
	reqBodyMap["contents"] = messagesPrepare(ignoreSystemPrompt, messages)
	marshal, _ := json.Marshal(reqBodyMap)

	req, _ := http.NewRequest("POST", "https://qianwen.biz.aliyun.com/dialog/conversation", bytes.NewReader(marshal))
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	//3. 发起请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	c.Writer.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	sessionId := ""

	msgChan := make(chan string, 10)

	go func() {
		defer close(msgChan)
		lastMsg := ""
		id := ""
		for {
			buf := make([]byte, 8192)
			n, bodyErr := resp.Body.Read(buf)
			if bodyErr != nil {
				tChunk := model.NewStopChatCompletionChunk(id, openaiBody.Model)
				tMarshal, _ := json.Marshal(tChunk)
				msgChan <- fmt.Sprintf("data: %s\n\n", tMarshal)
				msgChan <- fmt.Sprintf("data: [DONE]\n")
				break
			}
			var str = string(buf[:n])
			//截取{之后的内容
			index := strings.Index(str, "{")
			if index > 0 {
				str = str[index:]
			}
			var tongYiWebRespBody TongYiWebRespBody
			_ = json.Unmarshal([]byte(str), &tongYiWebRespBody)
			if sessionId == "" {
				sessionId = tongYiWebRespBody.SessionID
			}
			if len(tongYiWebRespBody.Contents) < 1 {
				continue
			}
			tMsg := tongYiWebRespBody.Contents[0].Content
			tMsg2 := strings.Replace(tMsg, lastMsg, "", 1)
			if tMsg2 == "" && tongYiWebRespBody.StopReason != "stop" {
				continue
			}
			if id == "" {
				id = tongYiWebRespBody.MsgID
			}

			chunk := model.NewChatCompletionChunk(id, tMsg2, "qwen-web")
			chunkBytes, _ := json.Marshal(chunk)

			msgChan <- fmt.Sprintf("data: %s\n\n", chunkBytes)
			lastMsg = tMsg
		}
	}()

	for msg := range msgChan {
		_, _ = c.Writer.WriteString(msg)
		c.Writer.Flush()
	}
	go func() {
		time.Sleep(time.Second * 1)
		deleteChat(key, sessionId)
	}()

}

//4. 处理返回结果

type TYMessage struct {
	Role        string `json:"role"`
	Content     any    `json:"content"`
	ContentType string `json:"contentType"`
}

func deleteChat(key, sessionId string) {
	//1. 准备发起请求的请求头
	var headers = fakeHeader
	cookieStr := generateCookie(key)
	headers["Cookie"] = cookieStr
	//2. 准备发起请求的请求体

	reqBodyMap := make(map[string]interface{})
	reqBodyMap["sessionId"] = sessionId
	marshal, _ := json.Marshal(reqBodyMap)

	req, _ := http.NewRequest("POST", "https://qianwen.biz.aliyun.com/dialog/session/delete", bytes.NewReader(marshal))
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	//3. 发起请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
}

var fakeHeader = map[string]string{
	"Accept":             "text/event-stream",
	"Accept-Encoding":    "gzip, deflate, br, zstd",
	"Content-Type":       "application/json",
	"Accept-Language":    "zh-CN,zh;q=0.9",
	"Cache-Control":      "no-cahce",
	"Origin":             "https://tongyi.aliyun.com",
	"Pragma":             "no-cache",
	"Sec-Ch-Ua":          "Chromium;v=122, Not(A:Brand;v=24, Google Chrome;v=122",
	"Sec-Ch-Ua-Mobile":   "?0",
	"Sec-Ch-Ua-Platform": "Windows",
	"Sec-Fetch-Dest":     "empty",
	"Sec-Fetch-Mode":     "cors",
	"Sec-Fetch-Site":     "same-site",
	"Referer":            "https://tongyi.aliyun.com/",
	"User-Agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36",
	"X-Platform":         "pc_tongyi",
	"X-Xsrf-Token":       "4506064f-1da8-4d22-a004-b36092db2789",
}

func messagesPrepare(ignoreSystemPrompt bool, messages []model.Message) []TYMessage {
	var resultContent strings.Builder
	resultContent.WriteString("<|im_start|>\n")
	for _, message := range messages {
		if ignoreSystemPrompt && message.Role == "system" {
			continue
		}
		resultContent.WriteString(message.Role)
		resultContent.WriteString(":")
		resultContent.WriteString(message.Content)
		resultContent.WriteString("\n")
	}
	resultContent.WriteString("<|im_end|>\n")
	return []TYMessage{
		{
			Role:        "user",
			ContentType: "text",
			Content:     resultContent.String(),
		},
	}
}

func generateCookie(key string) (cookieStr string) {
	arrs := []string{
		"login_tongyi_ticket=" + key,
		"_samesite_flag_=true",
		"t=" + uuid.NewString(),
		"channel=oug71n2fX3Jd5ualEfKACRvnsceUtpjUC5jHBpfWnSOXKhkvBNuSO8bG3v4HHjCgB722h7LqbHkB6sAxf3OvgA%3D%3D",
		"currentRegionId=cn-shenzhen",
		"aliyun_country=CN",
		"aliyun_lang=zh",
		"aliyun_site=CN",
	}
	cookieStr = strings.Join(arrs, "; ")
	return
}
