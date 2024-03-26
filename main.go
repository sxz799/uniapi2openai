package main

import (
	"github/sxz799/gemini2chatgpt/model"
	"github/sxz799/gemini2chatgpt/service"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
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
}

func main() {

	r := gin.Default()
	r.Use(Cors())
	r.GET("/", func(context *gin.Context) {
		context.String(200, "部署成功！[https://github.com/sxz799/gemini2chatgpt]")
	})
	r.POST("v1/chat/completions", func(c *gin.Context) {
		var apiKey string
		envApi := os.Getenv("API_KEY")
		if envApi == "" {
			auth := c.GetHeader("Authorization")
			if len(strings.Split(auth, " ")) != 2 {
				c.JSON(400, gin.H{
					"error": "Authorization header is invalid",
				})
				return
			}
			apiKey = strings.Split(auth, " ")[1]
		}

		var originBody model.ChatGPTRequestBody
		err := c.BindJSON(&originBody)
		if err != nil {
			c.JSON(400, gin.H{
				"error": "Request body is invalid",
			})
			return
		}
		service.DoTrans(apiKey, originBody, c)
	})
	r.NoRoute(func(c *gin.Context) {
		c.JSON(200, gin.H{
			"msg": "仅代理了`v1/chat/completions`接口",
		})
	})
	err := r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
