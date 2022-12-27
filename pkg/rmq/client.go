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
}

type Client struct {
	connectionUrl string
	kind          string
	durable       bool
	autoDelete    bool
	internal      bool
	noWait        bool
	exclusive     bool

	heartBeat time.Duration

	connection *amqp.Connection
	logger     *logrus.Logger
	broker     *Broker
}

func New(connectionUrl string, logger *logrus.Logger) *Client {
	return &Client{
		connectionUrl: connectionUrl,
		kind:          "direct",
		durable:       true,
		autoDelete:    false,
		internal:      false,
		noWait:        false,
		exclusive:     false,

		heartBeat: 60 * time.Second,

		logger: logger,
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

func (c *Client) CreateChannel(queueName, exchangeName, routingKey string) {
	c.logger.Info("trying to create a broker")
	channel, err := c.connection.Channel()
	if err != nil {
		c.logger.Fatalf("failed to create broker")
	}
	c.logger.Info("broker created: OK")
	c.broker = &Broker{
		channel:      channel,
		queueName:    queueName,
		exchangeName: exchangeName,
		routingKey:   routingKey,
	}
}

func (c *Client) DeclareExchange() {
	if err := c.broker.channel.ExchangeDeclare(
		c.broker.exchangeName,
		c.kind,
		c.durable,
		c.autoDelete,
		c.internal,
		c.noWait,
		nil,
	); err != nil {
		c.logger.Fatalf("failed to create exchange: '%s'", c.broker.exchangeName)
	}
}

func (c *Client) BindQueue() {
	if err := c.broker.channel.QueueBind(
		c.broker.queueName,
		c.broker.routingKey,
		c.broker.exchangeName,
		c.noWait,
		nil,
	); err != nil {
		c.logger.Fatalf("failed to bind queue and exchange: '%s'", c.broker.queueName)
	}
}

func (c *Client) CloseConnection() {
	err := c.connection.Close()
	if err != nil {
		c.logger.Errorf("failed to close connection")
	}
}

func (c *Client) CloseChannel() {
	err := c.broker.channel.Close()
	if err != nil {
		c.logger.Errorf("failed to close broker")
	}
}
