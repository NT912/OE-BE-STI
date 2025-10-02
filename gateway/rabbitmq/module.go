package rabbitmq

import (
	"log"

	"github.com/streadway/amqp"
)

type RabbitMQModule struct {
	Conn      *amqp.Connection
	Channel   *amqp.Channel
	RPCClient *RPCClient
}

var rabbitMQModule *RabbitMQModule

func InitModule(rabbitMQURL string) *RabbitMQModule {
	if rabbitMQModule != nil {
		return rabbitMQModule
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

	rpcClient, err := newRPCClient(ch)
	if err != nil {
		log.Fatalf("Failed to init RPC client: %v", err)
	}

	log.Println("✅ RabbitMQ and RPC Client initialized")

	rabbitMQModule = &RabbitMQModule{
		Conn:      conn,
		Channel:   ch,
		RPCClient: rpcClient,
	}
	return rabbitMQModule
}

func (m *RabbitMQModule) Publish(exchange, routingKey string, body []byte) error {
	return m.Channel.Publish(
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
