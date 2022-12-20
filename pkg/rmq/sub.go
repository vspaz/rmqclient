package rmq

import (
	"github.com/streadway/amqp"
	"os"
)

func (r *RmqClient) Consume(channel *amqp.Channel, consumerName string) <-chan amqp.Delivery {
	consumerChannel, err := channel.Consume(
		r.queueName,
		consumerName,
		true,
		r.exclusive,
		false,
		r.noWait,
		nil,
	)
	if err != nil {
		r.logger.Error("failed to consume")
		os.Exit(-1)
	}
	return consumerChannel
}

func (r *RmqClient) DeclareQueue(channel *amqp.Channel) {
	if _, err := channel.QueueDeclare(
		r.queueName,
		r.durable,
		r.autoDelete,
		r.exclusive,
		r.noWait,
		nil,
	); err != nil {
		r.logger.Error("failed to declare queue: ", r.queueName)
		os.Exit(-1)
	}
}
