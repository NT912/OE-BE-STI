package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	routes "gateway/api"
	"gateway/configs"
	"gateway/rabbitmq"
)

func init() {
	configs.InitEnv()
	configs.ConnectDatabase()
}

func main() {
	rabbitMQModule := rabbitmq.InitModule(configs.Env.RabbitMQURL)
	defer rabbitMQModule.Conn.Close()
	defer rabbitMQModule.Channel.Close()

	routersInit := routes.InitRouter(rabbitMQModule)
	port := configs.Env.Port

	endPoint := fmt.Sprintf("%s:%s", "0.0.0.0", port)
	maxHeaderBytes := 1 << 20

	server := &http.Server{
		Addr:           endPoint,
		Handler:        routersInit,
		MaxHeaderBytes: maxHeaderBytes,
	}

	go func() {
		log.Printf("✅ Start http server listening at: %s\n", endPoint)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Printf("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}
	log.Println("Server exiting")
}
