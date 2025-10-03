package rabbitmq

import (
	"log"
	"sync"

	"gateway/configs"

	"github.com/streadway/amqp"
)

var (
	RabbitSvc *RabbitMQService
	once      sync.Once
)

func InitModule() {
	once.Do(func() {
		rabbitMQURL := configs.Env.RabbitMQURL
		if rabbitMQURL == "" {
			log.Fatalf("❌ RABBITMQ_URL is not configured")
		}
		log.Printf("Connecting to RabbitMQ at: %s", rabbitMQURL)
		conn, err := amqp.Dial(rabbitMQURL)
		if err != nil {
			log.Fatalf("❌ Failed to connect RabbitMQ: %v", err)
		}
		ch, err := conn.Channel()
		if err != nil {
			log.Fatalf("❌ Failed to open RabbitMQ channel: %v", err)
		}

		svc, err := newRabbitMQService(conn, ch)
		if err != nil {
			log.Fatalf("❌ Failed to init RabbitMQ Service: %v", err)
		}

		RabbitSvc = svc
		log.Printf("✅ RabbitMQ Service initialized")
	})
}
