package kafka

import (
	"context"
	"errors"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type MessageHandler interface {
	Handle(ctx context.Context, msg kafka.Message) error
}

type Consumer struct {
	reader  *kafka.Reader
	handler MessageHandler
	logger  *zap.Logger
}

type Config struct {
	Brokers []string
	Topic   string
	GroupID string
	TimeOut time.Duration
}

func NewConsumer(cfg Config, handler MessageHandler, logger *zap.Logger) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:          cfg.Brokers,
		Topic:            cfg.Topic,
		GroupID:          cfg.GroupID,
		MinBytes:         10 << 13,
		MaxBytes:         10 << 23,
		RebalanceTimeout: cfg.TimeOut,
	})

	return &Consumer{
		reader:  reader,
		handler: handler,
		logger:  logger,
	}
}

func (c *Consumer) Run(ctx context.Context) {
	c.logger.Info("Consumer has been started")
	for {
		message, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				c.logger.Info("Context cancled or exceeded")
				break
			}
			c.logger.Error("Failed to read message", zap.Error(err))
			continue
		}
		c.logger.Info("Message recieved", zap.ByteString("value", message.Value))
		if err := c.handler.Handle(ctx, message); err != nil {
			// TODO: Make DLQ for failed messages
			c.logger.Error("Failed to handler message")
		}
		if err := c.reader.CommitMessages(ctx, message); err != nil {
			c.logger.Error("Failed to commit message", zap.Error(err))
		}
	}
}

func (c *Consumer) Close() error {
	c.logger.Info("Closing reader")
	return c.reader.Close()
}
