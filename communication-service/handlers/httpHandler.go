package handlers

import (
	"log"
	"net/http"

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
	var req services.ChatRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	log.Printf("Received request from gateway with prompt: %s", req.Prompt)

	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")

	err := h.aiService.GetAIStream(req.Prompt, ctx.Writer)
	if err != nil {
		log.Printf("Error from AiService stream: %v", err)
	}
}
