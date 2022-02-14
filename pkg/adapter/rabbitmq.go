package adapter

import (
	"context"

	"github.com/ivorscott/employee-service/pkg/msg"

	"github.com/ivorscott/employee-service/pkg/web"
	"github.com/wagslane/go-rabbitmq"
	"go.uber.org/zap"
)

var defaultPublishOptions = []func(options *rabbitmq.PublishOptions){
	rabbitmq.WithPublishOptionsContentType("application/json"),
	rabbitmq.WithPublishOptionsMandatory,
	rabbitmq.WithPublishOptionsPersistentDelivery,
	rabbitmq.WithPublishOptionsExchange("events"),
}

// Publish wraps the publishing logic and handles errors.
func (r *RabbitMQPublishAdapter) Publish(message []byte, routingKeys []string, options ...func(options *rabbitmq.PublishOptions)) {
	err := r.Publisher.Publish(
		message,
		routingKeys,
		append(options, defaultPublishOptions...)...,
	)
	if err != nil {
		r.log.Fatal("fatal error", zap.Error(err))
	}
}

type listenHandlerFunc func(ctx context.Context, message msg.Message) error

// Listen wraps the subscribing logic and handles errors.
func (r *RabbitMQListenAdapter) Listen(messageType msg.MessageType, queue string, routingKeys []string, handler listenHandlerFunc, options ...func(options *rabbitmq.ConsumeOptions)) {
	err := r.StartConsuming(
		func(d rabbitmq.Delivery) rabbitmq.Action {
			var err error

			message, err := msg.UnmarshalMessage(d.Body)
			if err != nil {
				r.log.Fatal("error decoding message", zap.Error(err))
			}

			if message.Type == messageType {
				err = handler(context.Background(), message)
				r.log.Error("error handling message", zap.Error(err))
			}

			switch err.(type) {
			case nil:
				return rabbitmq.Ack
			case *web.Error:
				return rabbitmq.NackDiscard
			default:
				return rabbitmq.NackRequeue
			}
		},
		queue,
		routingKeys,
		options...,
	)
	if err != nil {
		r.log.Fatal("fatal error", zap.Error(err))
	}
}
