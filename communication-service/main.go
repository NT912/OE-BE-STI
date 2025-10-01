package main

import (
	"log"
	"os"
	"os/signal"

	"communication-service/configs"
	handler "communication-service/handlers"
	"communication-service/services"
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

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("👋 Shutting down communication service...")
}
