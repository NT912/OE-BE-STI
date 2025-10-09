package ai

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

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
		aiRoutes.GET("/chat/ws", c.chatWebSocketHandler)
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

func (c *AiController) chatWebSocketHandler(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection to WebSocket: %v", err)
		return
	}
	defer conn.Close()
	log.Println("Client connected via WebSocket")

	// Test connect to WebSocket
	conn.WriteMessage(websocket.TextMessage, []byte("WebSocket connection established!"))

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}
		if messageType != websocket.TextMessage {
			continue
		}
		log.Printf("Received message from client: %s", string(p))

		var req ChatRequest
		if err := json.Unmarshal(p, &req); err != nil {
			log.Printf("Invalid JSON from client: %v", err)
			conn.WriteMessage(websocket.TextMessage, []byte("Error invalid JSON format."))
			continue

		}

		socketWriter := &SocketWriter{Conn: conn}
		if err := c.service.ChatStream(req, socketWriter); err != nil {
			log.Printf("Service stream error: %v", err)
			conn.WriteMessage(websocket.TextMessage, []byte("Error: Could not process stream."))
		}
	}
}
