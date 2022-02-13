// Package adapter creates new interfaces for outside services, like message stores and RPCs via HTTP clients.
package adapter

import (
	"os"

	"github.com/rabbitmq/amqp091-go"
	"github.com/wagslane/go-rabbitmq"
	"go.uber.org/zap"
)

// RabbitMQPublishAdapter is an adapter implementation for a RabbitMQ Publisher.
type RabbitMQPublishAdapter struct {
	log *zap.Logger
	*rabbitmq.Publisher
}

// RabbitMQListenAdapter is an adapter implementation for a RabbitMQ Listener.
type RabbitMQListenAdapter struct {
	log *zap.Logger
	rabbitmq.Consumer
}

// NewRabbitMQPublisher returns a new RabbitMQ publisher instance.
func NewRabbitMQPublisher(log *zap.Logger, url string, config amqp091.Config, optionFuncs ...func(*rabbitmq.PublisherOptions)) *RabbitMQPublishAdapter {
	publisher, err := rabbitmq.NewPublisher(url, config, optionFuncs...)
	if err != nil {
		log.Error("Fatal error", zap.Error(err))
		os.Exit(1)
	}
	return &RabbitMQPublishAdapter{log, publisher}
}

// NewRabbitMQListener returns a new  listener instance.
func NewRabbitMQListener(log *zap.Logger, url string, config amqp091.Config, optionFuncs ...func(options *rabbitmq.ConsumerOptions)) *RabbitMQListenAdapter {
	consumer, err := rabbitmq.NewConsumer(url, config, optionFuncs...)
	if err != nil {
		log.Error("Fatal error", zap.Error(err))
		os.Exit(1)
	}
	return &RabbitMQListenAdapter{log, consumer}
}
