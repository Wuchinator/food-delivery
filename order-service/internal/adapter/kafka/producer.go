package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Producer struct {
	writer *kafka.Writer
	logger *zap.Logger
}

type OrderCreatedEvent struct {
	OrderID   int64     `json:"order_id"`
	UserID    int64     `json:"user_id"`
	Timestamp time.Time `json:"timestamp"`
}

type Config struct {
	Brokers         []string
	Topic           string
	ProducerTimeout time.Duration
	RequireAcks     int
}

func NewProducer(cfg Config, logger *zap.Logger) *Producer {

	writer := &kafka.Writer{
		Addr:         kafka.TCP(cfg.Brokers...),
		Topic:        cfg.Topic,
		Balancer:     &kafka.CRC32Balancer{},
		WriteTimeout: cfg.ProducerTimeout,
		RequiredAcks: kafka.RequiredAcks(cfg.RequireAcks),
	}

	return &Producer{
		writer: writer,
		logger: logger.Named("kafka_producer"),
	}

}

func (p *Producer) SentOrCreated(ctx context.Context, event OrderCreatedEvent) error {
	valueBytes, err := json.Marshal(event)
	if err != nil {
		p.logger.Error("Failed to marshal event")
		return err
	}

	msg := kafka.Message{
		Key:   []byte(fmt.Sprintf("%d", event.UserID)),
		Value: valueBytes,
	}

	err = p.writer.WriteMessages(ctx, msg)
	if err != nil {
		p.logger.Error("Failed to write message", zap.Error(err))
		return err
	}

	p.logger.Info("OrderCreatedEvent sent to Kafka",
		zap.Int64("order_id", event.OrderID),
	)

	return nil
}

func (p *Producer) Close() error {
	p.logger.Info("Producer close")
	return p.writer.Close()
}
