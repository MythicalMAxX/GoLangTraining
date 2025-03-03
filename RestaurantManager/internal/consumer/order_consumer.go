package consumer

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"restaurantservice/pkg/delivery"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

type OrderItem struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

type Order struct {
	ID              string      `json:"id"`
	CustomerID      string      `json:"customer_id"`
	Items           []OrderItem `json:"items"`
	TotalAmount     float64     `json:"total_amount"`
	Status          string      `json:"status"`
	DeliveryAddress string      `json:"delivery_address"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}

type OrderMessage struct {
	Type      string    `json:"type"`
	Order     Order     `json:"order"`
	Timestamp time.Time `json:"timestamp"`
}

type OrderConsumer struct {
	channel        *amqp.Channel
	queue          string
	exchange       string
	routingKey     string
	restaurantID   uuid.UUID
	deliveryClient *delivery.DeliveryClient
}

func NewOrderConsumer(ch *amqp.Channel, queueName, exchange, routingKey string, restaurantID uuid.UUID) *OrderConsumer {
	return &OrderConsumer{
		channel:        ch,
		queue:          queueName,
		exchange:       exchange,
		routingKey:     routingKey,
		restaurantID:   restaurantID,
		deliveryClient: delivery.NewDeliveryClient(8086),
	}
}

func (c *OrderConsumer) StartConsumer() error {
	msgs, err := c.channel.Consume(
		c.queue, // queue
		"",      // consumer
		false,   // auto-ack
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)
	if err != nil {
		return err
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var orderMsg OrderMessage
			if err := json.Unmarshal(d.Body, &orderMsg); err != nil {
				log.Printf("Error parsing message: %v", err)
				d.Reject(false)
				continue
			}

			log.Printf("\nReceived Message:")
			log.Printf("Type: %s", orderMsg.Type)
			log.Printf("Order ID: %s", orderMsg.Order.ID)
			log.Printf("Customer ID: %s", orderMsg.Order.CustomerID)
			log.Printf("Status: %s", orderMsg.Order.Status)
			log.Printf("Total Amount: $%.2f", orderMsg.Order.TotalAmount)
			log.Printf("Items:")
			for _, item := range orderMsg.Order.Items {
				log.Printf("  - %s (Qty: %d, Price: $%.2f)", item.Name, item.Quantity, item.Price)
			}
			log.Printf("Timestamp: %v\n", orderMsg.Timestamp)

			switch orderMsg.Type {
			case "order.new":
				c.handleNewOrder(orderMsg)
			case "order.update":
				c.handleOrderUpdate(orderMsg)
			}

			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

	return nil
}

func (c *OrderConsumer) handleNewOrder(msg OrderMessage) {
	log.Printf("Processing new order: %s", msg.Order.ID)

	// Extract item names
	itemNames := make([]string, len(msg.Order.Items))
	for i, item := range msg.Order.Items {
		itemNames[i] = fmt.Sprintf("%s (x%d)", item.Name, item.Quantity)
	}

	// Create delivery request with proper format
	deliveryReq := &delivery.DeliveryRequest{
		OrderID:         msg.Order.ID, // Changed from ID to OrderID
		CustomerID:      msg.Order.CustomerID,
		RestaurantID:    c.restaurantID.String(),
		Status:          "PENDING",
		TotalAmount:     msg.Order.TotalAmount,
		DeliveryAddress: msg.Order.DeliveryAddress,
		Items:           itemNames, // Added items
	}

	// Log the request before sending
	log.Printf("Preparing delivery request: %+v", deliveryReq)

	// Send to delivery service
	err := c.deliveryClient.RequestDelivery(deliveryReq)
	if err != nil {
		log.Printf("Failed to request delivery: %v", err)
		msg.Order.Status = "DELIVERY_REQUEST_FAILED"
	} else {
		msg.Order.Status = "DELIVERY_REQUESTED"
		log.Printf("Successfully requested delivery for order: %s", msg.Order.ID)
	}

	// Publish status update back to message queue
	if err := c.publishOrderStatus(msg.Order); err != nil {
		log.Printf("Failed to publish order status update: %v", err)
	}
}

func (c *OrderConsumer) handleOrderUpdate(msg OrderMessage) {
	log.Printf("Handling order update for: %s, new status: %s",
		msg.Order.ID, msg.Order.Status)
}

func (c *OrderConsumer) publishOrderStatus(order Order) error {
	msg := OrderMessage{
		Type:      "order.update",
		Order:     order,
		Timestamp: time.Now(),
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return c.channel.Publish(
		c.exchange,     // exchange
		"order.update", // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)
}

// Helper function to extract item names
func getItemNames(items []OrderItem) []string {
	names := make([]string, len(items))
	for i, item := range items {
		names[i] = item.Name
	}
	return names
}
