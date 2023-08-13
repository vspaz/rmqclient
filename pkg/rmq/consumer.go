package rmq

import (
	"github.com/rabbitmq/amqp091-go"
	"os"
)

func (c *Channel) Consume(consumerName string) <-chan amqp091.Delivery {
	consumerChannel, err := c.channel.Consume(
		c.queueName,
		consumerName,
		true,
		c.exclusive,
		false,
		c.noWait,
		nil,
	)
	if err != nil {
		c.logger.Error("failed to consume")
		os.Exit(-1)
	}
	return consumerChannel
}
