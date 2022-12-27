package rmq

import (
	"github.com/streadway/amqp"
	"os"
)

func (b *Broker) Consume(consumerName string) <-chan amqp.Delivery {
	consumerChannel, err := b.channel.Consume(
		b.queueName,
		consumerName,
		true,
		b.exclusive,
		false,
		b.noWait,
		nil,
	)
	if err != nil {
		b.logger.Error("failed to consume")
		os.Exit(-1)
	}
	return consumerChannel
}
