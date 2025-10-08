package ai

import "github.com/gin-gonic/gin"

type AiController struct {
	service *AiService
}

func NewAiController(s *AiService) *AiController {
	return &AiController{service: s}
}

func (c *AiController) RegisterRoutes(r *gin.RouterGroup) {
}
