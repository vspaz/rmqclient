package rmq

import (
	"github.com/rabbitmq/amqp091-go"
	"time"
)

func (c *Channel) Publish(task []byte, contentType string) error {
	msg := amqp091.Publishing{
		DeliveryMode: amqp091.Persistent,
		ContentType:  contentType,
		Timestamp:    time.Now(),
		Body:         task,
	}

	err := c.channel.Publish(
		c.exchangeName,
		c.routingKey,
		false,
		false,
		msg,
	)
	if err != nil {
		c.logger.Error("failed to publish message")
	}
	return err
}
