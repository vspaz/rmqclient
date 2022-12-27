package rmq

import (
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"time"
)

type RmqClient struct {
	connectionUrl string
	kind          string
	durable       bool
	autoDelete    bool
	internal      bool
	noWait        bool
	exclusive     bool

	exchangeName string
	queueName    string
	routingKey   string
	heartBeat    time.Duration

	logger *logrus.Logger
}

func New(connectionUrl string, queueName string, exchangeName string, routingKey string, logger *logrus.Logger) *RmqClient {
	return &RmqClient{
		connectionUrl: connectionUrl,
		kind:          "direct",
		durable:       true,
		autoDelete:    false,
		internal:      false,
		noWait:        false,
		exclusive:     false,

		exchangeName: exchangeName,
		queueName:    queueName,
		routingKey:   routingKey,
		heartBeat:    60 * time.Second,

		logger: logger,
	}
}

func (r *RmqClient) Connect() *amqp.Connection {
	r.logger.Debugf("connecting to rabbitmq '%s'", r.connectionUrl)
	connection, err := amqp.DialConfig(r.connectionUrl, amqp.Config{Heartbeat: time.Second * r.heartBeat})
	if err != nil {
		r.logger.Fatalf("failed to establish connection at '%s'", r.connectionUrl)
	}
	r.logger.Info("connection to rabbitmq established at: OK")
	return connection
}

func (r *RmqClient) CreateChannel(connection *amqp.Connection) *amqp.Channel {
	r.logger.Info("trying to create a channel")
	channel, err := connection.Channel()
	if err != nil {
		r.logger.Fatalf("failed to create channel")
	}
	r.logger.Info("channel created: OK")
	return channel
}

func (r *RmqClient) DeclareExchange(chanel *amqp.Channel) {
	if err := chanel.ExchangeDeclare(
		r.exchangeName,
		r.kind,
		r.durable,
		r.autoDelete,
		r.internal,
		r.noWait,
		nil,
	); err != nil {
		r.logger.Fatalf("failed to create exchange: '%s'", r.exchangeName)
	}
}

func (r *RmqClient) BindQueue(channel *amqp.Channel) {
	if err := channel.QueueBind(
		r.queueName,
		r.routingKey,
		r.exchangeName,
		r.noWait,
		nil,
	); err != nil {
		r.logger.Fatalf("failed to bind queue and exchange: '%s'", r.queueName)
	}
}

func (r *RmqClient) CloseConnection(connection *amqp.Connection) {
	err := connection.Close()
	if err != nil {
		r.logger.Errorf("failed to close connection")
	}
}

func (r *RmqClient) CloseChannel(channel *amqp.Channel) {
	err := channel.Close()
	if err != nil {
		r.logger.Errorf("failed to close channel")
	}
}
