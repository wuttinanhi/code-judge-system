package services

import (
	"context"
	"log"
)

type kafkaMockService struct {
	ctx context.Context
}

// CreateTopic implements KafkaService.
func (*kafkaMockService) CreateTopic(topic string, partitions int) error {
	return nil
}

// Consume implements KafkaService.
func (*kafkaMockService) Consume(topic string, groupID string) (chan string, chan error) {
	// do nothing
	return nil, nil
}

// IsTopicExist implements KafkaService.
func (*kafkaMockService) IsTopicExist(topic string) bool {
	// do nothing
	return true
}

// OverriddenHost implements KafkaService.
func (*kafkaMockService) OverriddenHost(host string) {
	// do nothing
}

// Produce implements KafkaService.
func (*kafkaMockService) Produce(topic string, message string) error {
	// do nothing
	log.Println("KafkaMockService: Produce", topic, message)
	return nil
}

func NewKafkaMockService() KafkaService {
	return &kafkaMockService{
		ctx: context.Background(),
	}
}
