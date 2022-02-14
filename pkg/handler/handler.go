package handler

import "github.com/wagslane/go-rabbitmq"

type rabbitmqAdapter interface {
	Publish(message []byte, routingKey []string, options ...func(options *rabbitmq.PublishOptions))
}
