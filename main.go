package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sxz799/uniapi2openai/u2oService"
	"log"
	"net/http"
	"os"
)

func main() {

	servePort := os.Getenv("SERVE_PORT")
	if servePort == "" {
		servePort = "8080"
	}
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
		context.String(200, "部署成功!")
	})
	r.POST("/v1/chat/completions", func(c *gin.Context) {
		u2oService.DoTrans(true, c)
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
