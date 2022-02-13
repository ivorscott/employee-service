package handler

import "github.com/wagslane/go-rabbitmq"

type rabbitMQPublisher interface {
	Publish(message []byte, routingKey []string, options ...func(options *rabbitmq.PublishOptions))
}
