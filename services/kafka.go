package services

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaService interface {
	Produce(topic string, message string) error
	Consume(topic string, groupID string) (chan string, chan error)
	IsTopicExist(topic string) bool
	OverriddenHost(host string)
}

type kafkaService struct {
	host string
	ctx  context.Context
}

func newKafkaWriter(host string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(host),
		Balancer: &kafka.LeastBytes{},
	}
}

func getKafkaReader(host, topic, groupID string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{host},
		GroupID: groupID,
		Topic:   topic,
	})
}

// Produce implements KafkaService.
func (s *kafkaService) Produce(topic string, message string) error {
	if s.host == "" {
		panic("KafkaService: host is empty")
	}

	writer := newKafkaWriter(s.host)
	defer writer.Close()

	msg := kafka.Message{
		Topic: topic,
		Value: []byte(message),
		Key:   []byte(fmt.Sprintf("%d", time.Now().UnixNano())),
	}

	return writer.WriteMessages(s.ctx, msg)
}

// Consume implements KafkaService.
func (s *kafkaService) Consume(topic string, groupID string) (chan string, chan error) {
	reader := getKafkaReader(s.host, topic, groupID)

	// make channel
	messageC := make(chan string)
	errorC := make(chan error)

	go func() {
		for {
			msg, err := reader.ReadMessage(s.ctx)
			if err != nil {
				errorC <- err
				continue
			}

			messageStr := string(msg.Value)
			messageC <- messageStr
		}
	}()

	return messageC, errorC
}

func (s *kafkaService) IsTopicExist(topic string) bool {
	conn, err := kafka.DialContext(s.ctx, "tcp", s.host)
	if err != nil {
		return false
	}
	defer conn.Close()

	partitions, err := conn.ReadPartitions(topic)
	if err != nil {
		return false
	}

	if len(partitions) == 0 {
		return false
	}

	return true
}

func (s *kafkaService) OverriddenHost(host string) {
	s.host = host
}

func NewKafkaService(host string) KafkaService {

	return &kafkaService{
		host: host,
		ctx:  context.Background(),
	}
}
