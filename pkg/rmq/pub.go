package rmq

import (
	"github.com/streadway/amqp"
	"time"
)

func (c *Channel) PublishTask(task []byte, contentType string) error {
	msg := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
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
