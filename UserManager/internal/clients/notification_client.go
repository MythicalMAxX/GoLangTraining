package clients

import (
	"encoding/json"
	"fmt"
	"log"
	"userservices/internal/models"

	"github.com/IBM/sarama"
)

type NotificationClient struct {
	producer sarama.SyncProducer
	topic    string
}

func NewNotificationClient() (*NotificationClient, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %v", err)
	}

	return &NotificationClient{
		producer: producer,
		topic:    "user.notification",
	}, nil
}

func (c *NotificationClient) Close() error {
	return c.producer.Close()
}

func (c *NotificationClient) SendNotification(notification *models.Notification) error {
	// Validate priority
	switch notification.Priority {
	case models.High, models.Medium, models.Low:
		// Valid priority
	default:
		notification.Priority = models.Low // Default to low priority
	}

	// Convert notification to JSON
	jsonData, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %v", err)
	}

	// Create message
	msg := &sarama.ProducerMessage{
		Topic: c.topic,
		Value: sarama.StringEncoder(jsonData),
	}

	// Send message
	partition, offset, err := c.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send notification: %v", err)
	}

	log.Printf("Notification sent: Topic: %s, Partition: %d, Offset: %d", c.topic, partition, offset)
	return nil
}

func (c *NotificationClient) NotifySuccess(service, operation string) {
	notification := &models.Notification{
		Service:  service,
		Message:  fmt.Sprintf("Successfully completed %s", operation),
		Priority: models.Low,
	}

	if err := c.SendNotification(notification); err != nil {
		log.Printf("Failed to send success notification: %v", err)
	}
}

func (c *NotificationClient) NotifyError(service, operation string, err error) {
	notification := &models.Notification{
		Service:  service,
		Message:  fmt.Sprintf("Error in %s: %v", operation, err),
		Priority: models.High,
	}

	if err := c.SendNotification(notification); err != nil {
		log.Printf("Failed to send error notification: %v", err)
	}
}
