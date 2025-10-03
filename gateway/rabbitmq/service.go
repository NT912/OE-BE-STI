package rabbitmq

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

type RabbitMQService struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	rpcClient *RPCClient
}

func newRabbitMQService(conn *amqp.Connection, ch *amqp.Channel) (*RabbitMQService, error) {
	rpcClient, err := newRPCClient(ch)
	if err != nil {
		return nil, err
	}

	return &RabbitMQService{
		conn:      conn,
		channel:   ch,
		rpcClient: rpcClient,
	}, nil
}

func (s *RabbitMQService) Call(queueName string, requestBody interface{}) ([]byte, error) {
	return s.rpcClient.Call(queueName, requestBody)
}

func (s *RabbitMQService) Publish(exchange, routingKey string, body interface{}) error {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		log.Printf("Error marshalling publish payload: %v", err)
		return err
	}

	return s.channel.Publish(
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        bodyBytes,
		},
	)
}

func (s *RabbitMQService) Close() {
	if s.channel != nil {
		s.channel.Close()
	}
	if s.conn != nil {
		s.conn.Close()
	}
}
