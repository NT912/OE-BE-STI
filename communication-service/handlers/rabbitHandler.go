package handlers

import (
	"encoding/json"
	"log"

	"communication-service/services"

	"github.com/streadway/amqp"
)

type RabbitHandler struct{}

func NewRabbitHandler() *RabbitHandler {
	return &RabbitHandler{}
}

func (h *RabbitHandler) HandlerMessage(d amqp.Delivery) {
	switch d.RoutingKey {
	case "user.registered":
		h.handleUserRegistered(d)
	default:
		log.Printf("Unknown routing key: %s", d.RoutingKey)
	}
}

type UserRegisteredPayload struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (h *RabbitHandler) handleUserRegistered(d amqp.Delivery) {
	log.Printf("Received a message for user.registered: %s", d.Body)

	var payload UserRegisteredPayload
	err := json.Unmarshal(d.Body, &payload)
	if err != nil {
		log.Printf("Error devoding JSON: %s", err)
		return
	}

	err = services.SendWelcomeEmail(payload.Email, payload.Name)
	if err != nil {
		log.Printf("Error sending welcome email: %s", err)
		return
	}
	log.Printf("Welcome email sent to %s", payload.Email)
}

type RenderRequest struct {
	Template string                 `json:"template"`
	Data     map[string]interface{} `json:"data"`
}

func (h *RabbitHandler) HandleRPCRequest(d amqp.Delivery, ch *amqp.Channel) {
	log.Printf("Received a RPC request: %s", d.Body)

	var req RenderRequest
	err := json.Unmarshal(d.Body, &req)
	if err != nil {
		log.Printf("Error decoding RPC request: %s", err)
		d.Nack(false, false)
		return
	}

	renderedBody, err := services.RenderTemplate(req.Template, req.Data)
	if err != nil {
		log.Printf("Error rendering template: %s", err)
		d.Nack(false, false)
		return
	}

	err = ch.Publish(
		"",
		d.ReplyTo,
		false,
		false,
		amqp.Publishing{
			ContentType:   "text/html",
			CorrelationId: d.CorrelationId,
			Body:          []byte(renderedBody),
		},
	)
	if err != nil {
		log.Printf("Failed to publish RPC reply: %s", err)
	} else {
		log.Printf("RPC reply sent for CorrelationId: %s", d.CorrelationId)
	}
	d.Ack(false)
}
