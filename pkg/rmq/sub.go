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
