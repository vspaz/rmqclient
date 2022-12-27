package rmq

import (
	"github.com/streadway/amqp"
	"time"
)

func (r *RmqClient) PublishTask(channel *amqp.Channel, task []byte, contentType string) error {
	msg := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  contentType,
		Timestamp:    time.Now(),
		Body:         task,
	}

	err := channel.Publish(
		r.exchangeName,
		r.routingKey,
		false,
		false,
		msg,
	)
	if err != nil {
		r.logger.Error("failed to publish message")
	}
	return err
}
