package main

import (
	"log"
	"os"
	"os/signal"

	"communication-service/configs"
	handler "communication-service/handlers"
	"communication-service/services"

	"github.com/gin-gonic/gin"
)

func init() {
	configs.InitEnv()
	configs.ConnectRabbitMQ()
}

func main() {
	log.Println("✅ Communication service started...")

	rabbitService := services.NewRabbitMQService()
	defer rabbitService.Close()

	rabbitHandler := handler.NewRabbitHandler()

	rabbitService.SetupPubSubConsume(rabbitHandler.HandlerMessage)
	rabbitService.SetupRPCConsumer(rabbitHandler.HandleRPCRequest)

	r := gin.Default()
	httpHandler := handler.NewHttpHandler()
	httpHandler.RegisterRoutes(r)

	go func() {
		log.Println("Starting HTTP server on port 8081...")
		if err := r.Run(":8081"); err != nil {
			log.Fatalf("Failed to run HTTP service: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("👋 Shutting down communication service...")
}
