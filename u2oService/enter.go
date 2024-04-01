package u2oService

import (
	"github.com/gin-gonic/gin"
	"github.com/sxz799/uniapi2openai/config"
	"github.com/sxz799/uniapi2openai/model"
	"github.com/sxz799/uniapi2openai/u2oProducts/u2o4gemini"
	"github.com/sxz799/uniapi2openai/u2oProducts/u2o4tongyi"
	"github.com/sxz799/uniapi2openai/u2oProducts/u2o4tongyiWeb"
)

func DoTrans(ignoreSystemPrompt bool, c *gin.Context) {
	var openaiBody model.OpenaiBody
	err := c.BindJSON(&openaiBody)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Request body is invalid",
		})
		return
	}
	originModel := openaiBody.Model
	product, ok := config.ModelMap[originModel]
	if !ok {
		c.JSON(400, gin.H{
			"error": "没有找到您的模型",
		})
		return
	}
	switch product {
	case "gemini":
		u2o4gemini.DoTrans(ignoreSystemPrompt, openaiBody, c)
	case "tongyi":
		u2o4tongyi.DoTrans(ignoreSystemPrompt, openaiBody, c)
	case "qwen-web":
		u2o4tongyiWeb.DoTrans(ignoreSystemPrompt, openaiBody, c)
	}
}
