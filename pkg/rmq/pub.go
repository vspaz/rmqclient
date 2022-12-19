package rmq

import (
	"github.com/streadway/amqp"
	"time"
)

func (r *RmqConfig) PublishTask(channel *amqp.Channel, task []byte) error {
	msg := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		ContentType:  "application/json",
		Body:         task,
	}

	err := channel.Publish(
		"",
		r.exchangeName,
		false,
		false,
		msg,
	)
	if err != nil {
		r.logger.Error("failed to publish message")
	}
	return err
}
