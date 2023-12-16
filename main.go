package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github/sxz799/gemini2chatgpt/model"
	"github/sxz799/gemini2chatgpt/service"
	"log"
	"strings"
)

func main() {

	r := gin.Default()
	r.Use(cors.Default())
	r.GET("/", func(context *gin.Context) {
		context.String(200, "部署成功！")
	})

	r.POST("v1/chat/completions", func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if len(strings.Split(auth, " ")) != 2 {
			c.JSON(400, gin.H{
				"error": "Authorization header is invalid",
			})
			return
		}
		apiKey := strings.Split(auth, " ")[1]
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

	err := r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
