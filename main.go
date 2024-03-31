package main

import (
	"github.com/sxz799/gemini2chatgpt/gemini2chatgpt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	apiKey := os.Getenv("API_KEY")
	servePort := os.Getenv("SERVE_PORT")
	if servePort == "" {
		servePort = "8080"
	}
	ingoreSystemPrompt := os.Getenv("INGORE_SYSTEM_PROMPT") == "YES" || os.Getenv("INGORE_SYSTEM_PROMPT") == "yes"
	log.Println("API_KEY:", apiKey)
	log.Println("INGORE_SYSTEM_PROMPT:", ingoreSystemPrompt)
	r := gin.Default()
	cors := func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "*")
		c.Header("Access-Control-Allow-Headers", "*")
		c.Header("Access-Control-Expose-Headers", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
	r.Use(cors)
	r.GET("/", func(context *gin.Context) {
		context.String(200, "部署成功！[https://github.com/sxz799/gemini2chatgpt]")
	})
	r.POST("v1/chat/completions", func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if len(strings.Split(auth, " ")) != 2 {
			if apiKey == "" {
				c.JSON(400, gin.H{
					"error": "Authorization header is invalid",
				})
				return
			}
		} else {
			apiKey = strings.Split(auth, " ")[1]
		}
		gemini2chatgpt.DoTrans(ingoreSystemPrompt, "https://gemini.sxz799.top", apiKey, c)
	})
	r.NoRoute(func(c *gin.Context) {
		c.JSON(200, gin.H{
			"msg": "仅代理了`v1/chat/completions`接口",
		})
	})
	err := r.Run(":" + servePort)
	if err != nil {
		log.Fatal(err)
	}
}
