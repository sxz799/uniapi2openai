package u2o_service

import (
	"github.com/gin-gonic/gin"
	"github.com/sxz799/gemini2chatgpt/u2o_config"
	"github.com/sxz799/gemini2chatgpt/u2o_model"
	"github.com/sxz799/gemini2chatgpt/u2o_utils/u2o4gemini"
	"github.com/sxz799/gemini2chatgpt/u2o_utils/u2o4tongyi"
	"log"
)

func DoTrans(ignoreSystemPrompt bool, c *gin.Context) {
	var openaiBody u2o_model.OpenaiBody
	err := c.BindJSON(&openaiBody)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Request body is invalid",
		})
		return
	}
	originModel := openaiBody.Model
	log.Println(originModel)
	product, ok := u2o_config.ModelMap[originModel]
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
	}
}
