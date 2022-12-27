package rmq

import (
	"github.com/streadway/amqp"
	"time"
)

func (c *Client) PublishTask(task []byte, contentType string) error {
	msg := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  contentType,
		Timestamp:    time.Now(),
		Body:         task,
	}

	err := c.broker.channel.Publish(
		c.broker.exchangeName,
		c.broker.routingKey,
		false,
		false,
		msg,
	)
	if err != nil {
		c.logger.Error("failed to publish message")
	}
	return err
}
