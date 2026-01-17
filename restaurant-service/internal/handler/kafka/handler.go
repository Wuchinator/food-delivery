package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type OrderEventHandler struct {
	logger *zap.Logger
}

func NewOrderEventStruct(logger *zap.Logger) *OrderEventHandler {
	return &OrderEventHandler{
		logger: logger,
	}
}

func (e *OrderEventHandler) HandleMessage(ctx context.Context, message kafka.Message) error {

	return nil
}
