package handlers

import (
	"log"

	"communication-service/services"

	"github.com/gin-gonic/gin"
)

type HttpHandler struct {
	aiService *services.AiService
}

func NewHttpHandler(aiService *services.AiService) *HttpHandler {
	return &HttpHandler{
		aiService: aiService,
	}
}

func (h *HttpHandler) RegisterRoutes(r *gin.Engine) {
	r.GET("health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"status": "ok",
		})
	})

	internalRoutes := r.Group("/internal")
	{
		internalRoutes.POST("/ai/stream", h.handleAiStream)
	}
}

func (h *HttpHandler) handleAiStream(ctx *gin.Context) {
	log.Println("Received request from gateway to stream AI response.")

	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-control", "no-cache")
	ctx.Header("Connection", "keep-alive")

	err := h.aiService.GetAIStream(ctx.Writer)
	if err != nil {
		log.Printf("Error during AI stream: %v", err)
	}
}
