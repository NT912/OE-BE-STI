package ai

import (
	"log"

	"github.com/gin-gonic/gin"
)

type AiController struct {
	service *AiService
}

func NewAiController(s *AiService) *AiController {
	return &AiController{service: s}
}

func (c *AiController) RegisterRoutes(r *gin.RouterGroup) {
	aiRoutes := r.Group("/ai")
	{
		aiRoutes.POST("/chat", c.chatHandler)
	}
}

func (c *AiController) chatHandler(ctx *gin.Context) {
	var req ChatRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")

	err := c.service.ChatStream(req, ctx.Writer)
	if err != nil {
		log.Printf("Error during AI chat stream: %v", err)
	}
}
