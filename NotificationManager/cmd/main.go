package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/IBM/sarama"
	"github.com/fatih/color"
	"github.com/gorilla/mux"
	"github.com/hashicorp/consul/api"
	"github.com/joho/godotenv"
)

type Notification struct {
	Service   string    `json:"service"`
	Message   string    `json:"message"`
	Priority  string    `json:"priority"`
	Timestamp time.Time `json:"timestamp"`
}

func displayNotification(n Notification) {
	timestamp := n.Timestamp.Format("2006-01-02 15:04:05")

	// Create colored output based on priority
	priorityColor := color.New(color.FgWhite)
	switch n.Priority {
	case "high":
		priorityColor = color.New(color.FgRed, color.Bold)
	case "medium":
		priorityColor = color.New(color.FgYellow)
	case "low":
		priorityColor = color.New(color.FgGreen)
	}

	fmt.Printf("\n=== New Notification ===\n")
	fmt.Printf("Time: %s\n", timestamp)
	fmt.Printf("From: %s\n", color.CyanString(n.Service))
	fmt.Printf("Priority: ")
	priorityColor.Printf("%s\n", n.Priority)
	fmt.Printf("Message: %s\n", n.Message)
	fmt.Printf("=====================\n")
}

func registerService() (*api.Client, error) {
	config := api.DefaultConfig()
	config.Address = "localhost:8500"

	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	port, _ := strconv.Atoi(os.Getenv("SERVICE_PORT"))
	registration := &api.AgentServiceRegistration{
		ID:   "notification-service",
		Name: "notification-service",
		Port: port,
		Check: &api.AgentServiceCheck{
			HTTP:     fmt.Sprintf("http://localhost:%d/health", port),
			Interval: "10s",
			Timeout:  "5s",
		},
	}

	err = client.Agent().ServiceRegister(registration)
	return client, err
}

func setupKafkaConsumer() (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	consumer, err := sarama.NewConsumer([]string{os.Getenv("KAFKA_BROKER")}, config)
	if err != nil {
		return nil, err
	}

	return consumer, nil
}

func consumeMessages(consumer sarama.Consumer) {
	topic := os.Getenv("KAFKA_TOPIC")
	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Printf("Failed to start consumer for topic %s: %v", topic, err)
		return
	}

	go func() {
		for message := range partitionConsumer.Messages() {
			var notification Notification
			if err := json.Unmarshal(message.Value, &notification); err != nil {
				log.Printf("Error unmarshaling message: %v", err)
				continue
			}
			notification.Timestamp = time.Now()
			displayNotification(notification)
		}
	}()
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Register with Consul
	consulClient, err := registerService()
	if err != nil {
		log.Fatal("Failed to register service:", err)
	}
	defer consulClient.Agent().ServiceDeregister("notification-service")

	// Setup Kafka consumer
	consumer, err := setupKafkaConsumer()
	if err != nil {
		log.Fatal("Failed to setup Kafka consumer:", err)
	}
	defer consumer.Close()
	consumeMessages(consumer)

	router := mux.NewRouter()

	// Add health check endpoint for Consul
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods("GET")

	router.HandleFunc("/notify", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var notification Notification
		notification.Timestamp = time.Now()

		if err := json.NewDecoder(r.Body).Decode(&notification); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		displayNotification(notification)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Notification received"))
	}).Methods("POST")

	port := os.Getenv("SERVICE_PORT")
	fmt.Printf("Notification Service starting on :%s\n", port)

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		if err := http.ListenAndServe(":"+port, router); err != nil {
			log.Fatal(err)
		}
	}()

	<-stop
	fmt.Println("\nShutting down service...")
}
