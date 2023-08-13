package rmq

import (
	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"os"
)

type Channel struct {
	channel      *amqp091.Channel
	queueName    string
	exchangeName string
	routingKey   string

	kind       string
	durable    bool
	autoDelete bool
	internal   bool
	noWait     bool
	exclusive  bool

	logger *logrus.Logger
	conn   *Connection
}

func NewChannel(conn *Connection, queueName, exchangeName, routingKey string) *Channel {
	return &Channel{
		queueName:    queueName,
		exchangeName: exchangeName,
		routingKey:   routingKey,

		autoDelete: false,
		internal:   false,
		noWait:     false,
		exclusive:  false,
		durable:    true,

		logger: conn.logger,
		conn:   conn,
	}
}

func (c *Channel) Create() {
	c.logger.Info("trying to create a broker")
	channel, err := c.conn.connection.Channel()
	if err != nil {
		c.logger.Fatalf("failed to create broker")
	}
	c.logger.Info("broker created: OK")
	c.channel = channel
}

func (c *Channel) DeclareExchange(kind string) {
	if err := c.channel.ExchangeDeclare(
		c.exchangeName,
		kind,
		c.durable,
		c.autoDelete,
		c.internal,
		c.noWait,
		nil,
	); err != nil {
		c.logger.Fatalf("failed to create exchange: '%s'", c.exchangeName)
	}
}

func (c *Channel) DeclareQueue() {
	if _, err := c.channel.QueueDeclare(
		c.queueName,
		c.durable,
		c.autoDelete,
		c.exclusive,
		c.noWait,
		nil,
	); err != nil {
		c.logger.Error("failed to declare queue: ", c.queueName)
		os.Exit(-1)
	}
}

func (c *Channel) BindQueue() {
	if err := c.channel.QueueBind(
		c.queueName,
		c.routingKey,
		c.exchangeName,
		c.noWait,
		nil,
	); err != nil {
		c.logger.Fatalf("failed to bind queue and exchange: '%s'", c.queueName)
	}
}

func (c *Channel) Close() {
	err := c.channel.Close()
	if err != nil {
		c.logger.Errorf("failed to close broker")
	}
}
