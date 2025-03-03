package main

import (
	"encoding/json"
	"log"
	"notification-manager/internal/models"
	"os"
	"os/signal"
	"syscall"

	"github.com/streadway/amqp"
)

func main() {
	// Set up logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting Notification Service...")

	// Connect to RabbitMQ
	conn, err := amqp.Dial("amqp://admin:admin@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open channel: %v", err)
	}
	defer ch.Close()

	// Declare exchange
	err = ch.ExchangeDeclare(
		"orders_exchange",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare exchange: %v", err)
	}

	// Declare queue
	q, err := ch.QueueDeclare(
		"notification_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	// Bind queue to exchange
	err = ch.QueueBind(
		q.Name,
		"notification.#",
		"orders_exchange",
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to bind queue: %v", err)
	}

	// Start consuming messages
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
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	// Handle graceful shutdown
	done := make(chan bool)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for d := range msgs {
			var notification models.Notification
			if err := json.Unmarshal(d.Body, &notification); err != nil {
				log.Printf("Error unmarshaling notification: %v", err)
				continue
			}

			// Display notification in a formatted way
			log.Printf("\n=== New Notification ===\nType: %s\nOrder ID: %s\nMessage: %s\n",
				notification.Type,
				notification.OrderID,
				notification.Message,
			)
			if notification.DeliveryID != "" {
				log.Printf("Delivery ID: %s\n", notification.DeliveryID)
			}
			log.Println("=====================")
		}
	}()

	log.Println("🔔 Notification Service is running. Press CTRL+C to exit.")

	// Wait for shutdown signal
	<-sigChan
	log.Println("Shutting down notification service...")
	done <- true
}
