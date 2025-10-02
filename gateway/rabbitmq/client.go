package rabbitmq

import (
	"encoding/json"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

type RPCClient struct {
	ch      *amqp.Channel
	replies <-chan amqp.Delivery
	replyTo string
	corrMap map[string]chan []byte
	mu      sync.Mutex
}

func newRPCClient(ch *amqp.Channel) (*RPCClient, error) {
	if ch == nil {
		return nil, errors.New("RabbitMQ channel is not initialized")
	}

	q, err := ch.QueueDeclare(
		"",
		false,
		true,
		true,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	client := &RPCClient{
		ch:      ch,
		replies: msgs,
		replyTo: q.Name,
		corrMap: make(map[string]chan []byte),
	}
	go client.handleReplies()
	return client, nil
}

func (c *RPCClient) handleReplies() {
	for d := range c.replies {
		c.mu.Lock()
		ch, ok := c.corrMap[d.CorrelationId]
		c.mu.Unlock()
		if ok {
			ch <- d.Body
		}
	}
}

func (c *RPCClient) Call(rpcQueue string, requestBody interface{}) ([]byte, error) {
	corrID := uuid.New().String()
	callback := make(chan []byte)

	c.mu.Lock()
	c.corrMap[corrID] = callback
	c.mu.Unlock()

	defer func() {
		c.mu.Lock()
		delete(c.corrMap, corrID)
		c.mu.Unlock()
	}()

	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}
	err = c.ch.Publish(
		"",
		rpcQueue,
		false,
		false,
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: corrID,
			ReplyTo:       c.replyTo,
			Body:          body,
		})
	if err != nil {
		return nil, err
	}
	log.Printf("[x] Sent RPC request with CorrID: %s", corrID)
	select {
	case res := <-callback:
		log.Printf("[.] Got RPC response for CorrID: %s", corrID)
		return res, nil
	case <-time.After(5 * time.Second):
		return nil, errors.New("RPC request timed out")
	}
}
