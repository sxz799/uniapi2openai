# gemini2chatgpt

## 介绍
一个简单的小工具,将不同模型的api转为openai的api
以支持b端或c端gpt应用[[ChatGPT-Next-Web](https://github.com/ChatGPTNextWeb/ChatGPT-Next-Web)]


## 使用方法

将程序的自定义接口设置为你服务器的地址  

只映射了'/v1/chat/completions' 接口




## 交叉编译

```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o gemini2chatgpt
```



测试接口
```
curl --request POST \
  --url https://gemini2chatgpt-xxxxxx.onrender.com/v1/chat/completions \
  --header 'Authorization: Bearer YOUR_GEMINI_API' \
  --header 'content-type: application/json' \
  --data '{
  "model": "gpt-4",
  "messages": [
    {
      "role": "user",
      "content": "你好"
    }
  ],
  "stream": true
}'
```



