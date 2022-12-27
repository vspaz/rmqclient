package rmq

import (
	"github.com/streadway/amqp"
	"time"
)

func (b *Broker) PublishTask(task []byte, contentType string) error {
	msg := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  contentType,
		Timestamp:    time.Now(),
		Body:         task,
	}

	err := b.channel.Publish(
		b.exchangeName,
		b.routingKey,
		false,
		false,
		msg,
	)
	if err != nil {
		b.logger.Error("failed to publish message")
	}
	return err
}
