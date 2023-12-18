# gemini2chatgpt

## 介绍
一个简单的小工具,将gemini的api转为chatgpt的api  
方便让各种chatgpt的b端或c端gpt应用[[ChatGPT-Next-Web](https://github.com/ChatGPTNextWeb/ChatGPT-Next-Web)]方便的使用gemini模型  



## 使用方法

将程序部署到一台可以连接google的服务器上

将聊天程序的自定义接口设置为你服务器的地址  
apikey填写Google gemini的api

推荐使用[Render](https://dashboard.render.com/)部署

### Render部署方法
New-> Web Service -> Deploy an existing image from a registry -> Image URL填写`sxz799/gemini2chatgpt:latest`

部署完成后会给一个链接 类似 `https://gemini2chatgpt-xxxxxx.onrender.com`

测试方法
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

  

## 使用docker

```
docker run -d --restart always --name gemini2chatgpt -p 8080:8080 sxz799/gemini2chatgpt:latest
```

