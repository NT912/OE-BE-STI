package ai

import "github.com/gin-gonic/gin"

func InitModule(r *gin.Engine) {
	aiService := NewAiService()
	controller := NewAiController(aiService)

	api := r.Group("/api/v1")
	controller.RegisterRoutes(api)
}
