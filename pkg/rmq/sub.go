package rmq

import (
	"github.com/streadway/amqp"
	"os"
)

func (r *RmqClient) Consume(queueName, consumerName string) <-chan amqp.Delivery {
	consumerChannel, err := r.channel.Consume(
		queueName,
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

func (r *RmqClient) DeclareQueue(queueName string) {
	if _, err := r.channel.QueueDeclare(
		queueName,
		r.durable,
		r.autoDelete,
		r.exclusive,
		r.noWait,
		nil,
	); err != nil {
		r.logger.Error("failed to declare queue: ", queueName)
		os.Exit(-1)
	}
}
