# uniapi2openai

## 介绍
一个简单的小工具,将不同模型的api转为openai的api
以支持b端或c端gpt应用[[ChatGPT-Next-Web](https://github.com/ChatGPTNextWeb/ChatGPT-Next-Web)]

目前已接入 

* gemini 
* 阿里通义千问 
* 通义千问web转api

## 支持模型列表

-	"gemini-pro"
-	"qwen-turbo"
-	"qwen-1.8b-chat"
-	"qwen-web"

## 使用方法

将程序的自定义接口设置为你服务器的地址  

只映射了'/v1/chat/completions' 接口




## 交叉编译

```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o uniapi2openai
```



测试接口
```
curl --request POST \
  --url https://uniapi2openai-xxxxxx.onrender.com/v1/chat/completions \
  --header 'Authorization: Bearer YOUR_API_KEY' \
  --header 'content-type: application/json' \
  --data '{
  "model": "YOUR_MODEL",
  "messages": [
    {
      "role": "user",
      "content": "你好"
    }
  ],
  "stream": true
}'
```



