package messaging

import (
	"encoding/json"
	"time"

	"github.com/IBM/sarama"
)

const (
	TopicInvCreate      = "admin.inv.create"
	TopicInvCreateResp  = "admin.inv.create.response"
	TopicInvUpdate      = "admin.inv.update"
	TopicInvUpdateResp  = "admin.inv.update.response"
	TopicInvRead        = "user.inv.read"
	TopicInvReadResp    = "user.inv.read.response"
	TopicInvReadAll     = "user.inv.readall"
	TopicInvReadAllResp = "user.inv.readall.response"
)

type KafkaMessage struct {
	ID        string      `json:"id"`
	Data      interface{} `json:"data"`
	Status    string      `json:"status"`
	Error     string      `json:"error,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

type KafkaClient struct {
	producer sarama.SyncProducer
	consumer sarama.Consumer
}

func NewKafkaClient() (*KafkaClient, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, config)
	if err != nil {
		return nil, err
	}

	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, nil)
	if err != nil {
		return nil, err
	}

	return &KafkaClient{
		producer: producer,
		consumer: consumer,
	}, nil
}

func (k *KafkaClient) PublishMessage(topic string, messageID string, data interface{}, status string) error {
	message := KafkaMessage{
		ID:        messageID,
		Data:      data,
		Status:    status,
		Timestamp: time.Now(),
	}

	json, err := json.Marshal(message)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(json),
		Key:   sarama.StringEncoder(messageID),
	}

	_, _, err = k.producer.SendMessage(msg)
	return err
}

func (k *KafkaClient) ConsumeMessages(topic string, handler func(*KafkaMessage) error) error {
	partitionConsumer, err := k.consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		return err
	}

	go func() {
		for msg := range partitionConsumer.Messages() {
			var kafkaMsg KafkaMessage
			if err := json.Unmarshal(msg.Value, &kafkaMsg); err != nil {
				continue
			}

			if err := handler(&kafkaMsg); err != nil {
				// Publish error response to corresponding response topic
				k.PublishMessage(topic+"response", kafkaMsg.ID, nil, "ERROR")
			}
		}
	}()

	return nil
}
