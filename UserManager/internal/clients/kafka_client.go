package clients

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

type KafkaClient struct {
	producer         sarama.SyncProducer
	consumer         sarama.Consumer
	messageResponses map[string]chan []byte
	errorResponses   map[string]chan error
	mutex            sync.RWMutex
}

type AuthMessage struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	Action string `json:"action"` // "generate" or "validate"
	Token  string `json:"token,omitempty"`
}

type AuthResponse struct {
	Token     string `json:"token,omitempty"`
	IsValid   int    `json:"isValid,omitempty"`
	UUID      string `json:"uuid,omitempty"`
	Role      string `json:"role,omitempty"`
	Error     string `json:"error,omitempty"`
	Action    string `json:"action"`
	RequestID string `json:"request_id"`
}

func NewKafkaClient() (*KafkaClient, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Timeout = time.Second * 10
	config.Consumer.Fetch.Min = 1
	config.Consumer.Fetch.Default = 1024 * 1024
	config.Consumer.MaxWaitTime = 100 * time.Millisecond
	config.Consumer.Return.Errors = true

	// Connect to Kafka
	brokers := []string{"localhost:9092"}
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %v", err)
	}

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		producer.Close()
		return nil, fmt.Errorf("failed to create consumer: %v", err)
	}

	client := &KafkaClient{
		producer:         producer,
		consumer:         consumer,
		messageResponses: make(map[string]chan []byte),
		errorResponses:   make(map[string]chan error),
	}

	// Start response consumers
	if err := client.startResponseConsumers(); err != nil {
		producer.Close()
		consumer.Close()
		return nil, err
	}

	log.Println("Successfully connected to Kafka")
	return client, nil
}

func (c *KafkaClient) startResponseConsumers() error {
	// List of response topics to monitor
	responseTopics := []string{
		"user.auth.response",
		"admin.inv.create.response",
		"admin.inv.update.response",
		"user.inv.read.response",
		"user.inv.readall.response",
	}

	for _, topic := range responseTopics {
		if err := c.consumeResponses(topic); err != nil {
			return err
		}
	}
	return nil
}

func (c *KafkaClient) consumeResponses(topic string) error {
	partitionConsumer, err := c.consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		return fmt.Errorf("error creating consumer for %s: %v", topic, err)
	}

	go func() {
		for msg := range partitionConsumer.Messages() {
			var response struct {
				ID     string          `json:"id"`
				Data   json.RawMessage `json:"data"`
				Status string          `json:"status"`
				Error  string          `json:"error,omitempty"`
			}

			if err := json.Unmarshal(msg.Value, &response); err != nil {
				log.Printf("Error unmarshaling response from %s: %v", topic, err)
				continue
			}

			c.mutex.RLock()
			respChan, exists := c.messageResponses[response.ID]
			errChan, errExists := c.errorResponses[response.ID]
			c.mutex.RUnlock()

			if exists && errExists {
				if response.Error != "" {
					errChan <- fmt.Errorf(response.Error)
				} else {
					respChan <- response.Data
				}
			}
		}
	}()

	return nil
}

func (c *KafkaClient) Close() {
	if c.producer != nil {
		c.producer.Close()
	}
	if c.consumer != nil {
		c.consumer.Close()
	}
}

func (c *KafkaClient) SendAuthRequest(userID, role, action, token string) (*AuthResponse, error) {
	msg := AuthMessage{
		UserID: userID,
		Role:   role,
		Action: action,
		Token:  token,
	}

	jsonData, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	// Send message to user.auth topic
	_, _, err = c.producer.SendMessage(&sarama.ProducerMessage{
		Topic: "user.auth",
		Value: sarama.StringEncoder(jsonData),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %v", err)
	}

	log.Printf("Sent auth request: %s", string(jsonData))

	// Listen for response
	partitionConsumer, err := c.consumer.ConsumePartition("user.auth.response", 0, sarama.OffsetNewest)
	if err != nil {
		return nil, fmt.Errorf("failed to create partition consumer: %v", err)
	}
	defer partitionConsumer.Close()

	// Wait for response with timeout
	timeout := time.After(5 * time.Second)
	select {
	case msg := <-partitionConsumer.Messages():
		var response AuthResponse
		if err := json.Unmarshal(msg.Value, &response); err != nil {
			return nil, fmt.Errorf("failed to parse response: %v", err)
		}
		log.Printf("Received auth response: %+v", response)
		return &response, nil
	case <-timeout:
		return nil, fmt.Errorf("timeout waiting for response")
	}
}

func (c *KafkaClient) PublishAndWaitForResponse(topic string, message interface{}, responseTopic string) ([]byte, error) {
	messageID := uuid.New().String()

	// Create response channels
	respChan := make(chan []byte, 1)
	errChan := make(chan error, 1)

	// Register channels
	c.mutex.Lock()
	c.messageResponses[messageID] = respChan
	c.errorResponses[messageID] = errChan
	c.mutex.Unlock()

	// Cleanup when done
	defer func() {
		c.mutex.Lock()
		delete(c.messageResponses, messageID)
		delete(c.errorResponses, messageID)
		c.mutex.Unlock()
		close(respChan)
		close(errChan)
	}()

	// Create and send message
	kafkaMessage := struct {
		ID        string      `json:"id"`
		Data      interface{} `json:"data"`
		Timestamp time.Time   `json:"timestamp"`
	}{
		ID:        messageID,
		Data:      message,
		Timestamp: time.Now(),
	}

	jsonData, err := json.Marshal(kafkaMessage)
	if err != nil {
		return nil, fmt.Errorf("error marshaling message: %v", err)
	}

	_, _, err = c.producer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(messageID),
		Value: sarama.StringEncoder(jsonData),
	})
	if err != nil {
		return nil, fmt.Errorf("error sending message: %v", err)
	}

	log.Printf("Published message to %s with ID: %s", topic, messageID)

	// Wait for response
	select {
	case data := <-respChan:
		return data, nil
	case err := <-errChan:
		return nil, err
	case <-time.After(30 * time.Second):
		return nil, fmt.Errorf("timeout waiting for response from %s", responseTopic)
	}
}
