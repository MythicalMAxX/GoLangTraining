package messaging

import (
	"encoding/json"
	"log"
	"restaurantservice/internal/consumer"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

type RabbitMQClient struct {
	conn     *amqp.Connection
	ch       *amqp.Channel
	exchange string
}

const (
	OrdersExchange          = "orders_exchange"
	RestaurantExchange      = "restaurant_events"
	OrdersQueue             = "restaurant_orders_queue"
	RestaurantAssignedTopic = "restaurant.assigned"
	OrdersTopic             = "orders.new"
	OrderRoutingKey         = "order.#"
)

func NewRabbitMQClient(url string, restaurantID uuid.UUID) (*RabbitMQClient, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	// Declare the exchange
	err = ch.ExchangeDeclare(
		OrdersExchange, // name
		"topic",        // type
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		return nil, err
	}

	// Declare the queue
	q, err := ch.QueueDeclare(
		OrdersQueue, // name
		true,        // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		return nil, err
	}

	// Bind queue to exchange
	err = ch.QueueBind(
		q.Name,          // queue name
		OrderRoutingKey, // routing key
		OrdersExchange,  // exchange
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	// Start the consumer
	orderConsumer := consumer.NewOrderConsumer(ch, OrdersQueue, OrdersExchange, OrderRoutingKey, restaurantID)
	go func() {
		if err := orderConsumer.StartConsumer(); err != nil {
			log.Printf("Error starting consumer: %v", err)
		}
	}()

	client := &RabbitMQClient{
		conn:     conn,
		ch:       ch,
		exchange: OrdersExchange,
	}

	log.Println("Successfully connected to RabbitMQ and started consumer")
	return client, nil
}

func (c *RabbitMQClient) ConsumeOrders(handler func([]byte) error) error {
	msgs, err := c.ch.Consume(
		OrdersQueue, // queue
		"",          // consumer
		false,       // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			err := handler(msg.Body)
			if err != nil {
				log.Printf("Error processing message: %v", err)
				msg.Nack(false, true) // Negative acknowledgment, requeue the message
			} else {
				msg.Ack(false) // Acknowledge the message
			}
		}
	}()

	log.Printf("Started consuming messages from queue: %s", OrdersQueue)
	return nil
}

func (c *RabbitMQClient) PublishRestaurantAssigned(message interface{}) error {
	return c.PublishMessage(RestaurantAssignedTopic, message)
}

func (c *RabbitMQClient) PublishMessage(routingKey string, message interface{}) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = c.ch.Publish(
		c.exchange, // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         data,
		})

	if err != nil {
		log.Printf("Failed to publish message: %v", err)
		return err
	}

	log.Printf("Successfully published message to %s with routing key %s", c.exchange, routingKey)
	return nil
}

func (c *RabbitMQClient) Close() {
	if c.ch != nil {
		c.ch.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}
