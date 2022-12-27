package rmq

import (
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"time"
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

	logger *logrus.Logger
}

type Client struct {
	connectionUrl string
	heartBeat     time.Duration

	connection *amqp.Connection
	logger     *logrus.Logger
}

func New(connectionUrl string, logger *logrus.Logger) *Client {
	return &Client{
		connectionUrl: connectionUrl,
		heartBeat:     60 * time.Second,
		logger:        logger,
	}
}

func (c *Client) Connect() {
	c.logger.Debugf("connecting to rabbitmq '%s'", c.connectionUrl)
	connection, err := amqp.DialConfig(c.connectionUrl, amqp.Config{Heartbeat: time.Second * c.heartBeat})
	if err != nil {
		c.logger.Fatalf("failed to establish connection at '%s'", c.connectionUrl)
	}
	c.logger.Info("connection to rabbitmq established at: OK")
	c.connection = connection
}

func (c *Client) CreateChannel(queueName, exchangeName, routingKey string) *Broker {
	c.logger.Info("trying to create a broker")
	channel, err := c.connection.Channel()
	if err != nil {
		c.logger.Fatalf("failed to create broker")
	}
	c.logger.Info("broker created: OK")
	return &Broker{
		channel:      channel,
		queueName:    queueName,
		exchangeName: exchangeName,
		routingKey:   routingKey,

		kind:       "direct",
		durable:    true,
		autoDelete: false,
		internal:   false,
		noWait:     false,
		exclusive:  false,

		logger: c.logger,
	}
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
		b.logger.Fatalf("failed to create exchange: '%s'", c.broker.exchangeName)
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
		b.logger.Fatalf("failed to bind queue and exchange: '%s'", c.broker.queueName)
	}
}

func (c *Client) CloseConnection() {
	err := c.connection.Close()
	if err != nil {
		c.logger.Errorf("failed to close connection")
	}
}

func (b *Broker) CloseChannel() {
	err := b.channel.Close()
	if err != nil {
		b.logger.Errorf("failed to close broker")
	}
}
