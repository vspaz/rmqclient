package rmq

import (
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type Broker struct {
	channel      *amqp.Channel
	queueName    string
	exchangeName string
	routingKey   string

	kind       string
	durable    bool
	autoDelete bool
	internal   bool
	noWait     bool
	exclusive  bool

	logger     *logrus.Logger
	connection *amqp.Connection
}

func NewBroker(queueName, exchangeName, routingKey string, client *Client) *Broker {
	return &Broker{
		queueName:    queueName,
		exchangeName: exchangeName,
		routingKey:   routingKey,

		kind:       "direct",
		durable:    true,
		autoDelete: false,
		internal:   false,
		noWait:     false,
		exclusive:  false,

		logger:     client.logger,
		connection: client.connection,
	}
}

func (b *Broker) CreateChannel(connection *amqp.Connection) {
	b.logger.Info("trying to create a broker")
	channel, err := connection.Channel()
	if err != nil {
		b.logger.Fatalf("failed to create broker")
	}
	b.logger.Info("broker created: OK")
	b.channel = channel
}

func (b *Broker) DeclareExchange() {
	if err := b.channel.ExchangeDeclare(
		b.exchangeName,
		b.kind,
		b.durable,
		b.autoDelete,
		b.internal,
		b.noWait,
		nil,
	); err != nil {
		b.logger.Fatalf("failed to create exchange: '%s'", b.exchangeName)
	}
}

func (b *Broker) BindQueue() {
	if err := b.channel.QueueBind(
		b.queueName,
		b.routingKey,
		b.exchangeName,
		b.noWait,
		nil,
	); err != nil {
		b.logger.Fatalf("failed to bind queue and exchange: '%s'", b.queueName)
	}
}

func (b *Broker) CloseChannel() {
	err := b.channel.Close()
	if err != nil {
		b.logger.Errorf("failed to close broker")
	}
}
