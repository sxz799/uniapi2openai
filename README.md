# gemini2chatgpt

## 介绍
一个简单的小工具,将gemini的api转为chatgpt的api  
方便让各种chatgpt的b端或c端gpt应用[[ChatGPT-Next-Web](https://github.com/ChatGPTNextWeb/ChatGPT-Next-Web)]
方便使用gemini模型  


## 使用方法

将程序部署到一台可以连接google的服务器上

将聊天程序的自定义接口设置为你服务器的地址  

只映射了'/v1/chat/completions' 接口

apikey填写Google gemini的api

推荐使用[Render](https://dashboard.render.com/)部署




### Render部署方法

[详细教程](https://blog.sxz799.xyz/posts/%E6%8A%80%E5%B7%A7/2023-12-19%E8%AE%A9chatgpt%E5%AE%A2%E6%88%B7%E7%AB%AF%E7%94%A8%E4%B8%8Agoogle%E5%AE%B6%E7%9A%84gemini-pro%E6%A8%A1%E5%9E%8B/)  
New-> Web Service -> Deploy an existing image from a registry -> Image URL填写`sxz799/gemini2chatgpt:latest`

部署完成后会给一个链接 类似 `https://gemini2chatgpt-xxxxxx.onrender.com`



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

## 关于API_KEY

程序会优先在环境变量中读取`API_KEY`  
如果没有找到才会从请求头中读取 

  

## 使用docker

```
docker run -d --restart always --name gemini2chatgpt -p 8080:8080 sxz799/gemini2chatgpt:latest
```

