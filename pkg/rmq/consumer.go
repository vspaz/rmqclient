package rmq

import (
	"github.com/streadway/amqp"
	"os"
)

func (c *Channel) Consume(consumerName string) <-chan amqp.Delivery {
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
