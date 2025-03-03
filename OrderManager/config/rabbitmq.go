package config

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

const (
	OrdersExchange   = "orders_exchange"
	OrdersQueue      = "orders_queue"
	OrdersRoutingKey = "order.#"

	// Message types
	MessageTypeNewOrder    = "order.new"
	MessageTypeUpdateOrder = "order.update"
	MessageTypeDeleteOrder = "order.delete"
)

func ConnectToRabbitMQ(url string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}
	return conn, nil
}

func SetupRabbitMQ(ch *amqp.Channel) error {
	// Declare exchange
	err := ch.ExchangeDeclare(
		OrdersExchange, // name
		"topic",        // changed type to topic
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %v", err)
	}

	// Declare queue
	_, err = ch.QueueDeclare(
		OrdersQueue, // name
		true,        // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %v", err)
	}

	// Bind queue to exchange with topic routing key
	err = ch.QueueBind(
		OrdersQueue,      // queue name
		OrdersRoutingKey, // routing key pattern
		OrdersExchange,   // exchange
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to bind queue: %v", err)
	}

	log.Println("RabbitMQ setup completed successfully")
	return nil
}
