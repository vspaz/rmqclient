package rmq

import (
	"github.com/streadway/amqp"
	"os"
)

func (c *Client) Consume(queueName, consumerName string) <-chan amqp.Delivery {
	consumerChannel, err := c.broker.channel.Consume(
		queueName,
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

func (c *Client) DeclareQueue(queueName string) {
	if _, err := c.broker.channel.QueueDeclare(
		queueName,
		c.durable,
		c.autoDelete,
		c.exclusive,
		c.noWait,
		nil,
	); err != nil {
		c.logger.Error("failed to declare queue: ", queueName)
		os.Exit(-1)
	}
}
