package rmq

import (
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"time"
)

type Channel struct {
	channel      *amqp.Channel
	queueName    string
	exchangeName string
	routingKey   string
}

type RmqClient struct {
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
	channel    *Channel
}

func New(connectionUrl string, logger *logrus.Logger) *RmqClient {
	return &RmqClient{
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

func (r *RmqClient) Connect() {
	r.logger.Debugf("connecting to rabbitmq '%s'", r.connectionUrl)
	connection, err := amqp.DialConfig(r.connectionUrl, amqp.Config{Heartbeat: time.Second * r.heartBeat})
	if err != nil {
		r.logger.Fatalf("failed to establish connection at '%s'", r.connectionUrl)
	}
	r.logger.Info("connection to rabbitmq established at: OK")
	r.connection = connection
}

func (r *RmqClient) CreateChannel(queueName, exchangeName, routingKey string) {
	r.logger.Info("trying to create a channel")
	channel, err := r.connection.Channel()
	if err != nil {
		r.logger.Fatalf("failed to create channel")
	}
	r.logger.Info("channel created: OK")
	r.channel = &Channel{
		channel:      channel,
		queueName:    queueName,
		exchangeName: exchangeName,
		routingKey:   routingKey,
	}
}

func (r *RmqClient) DeclareExchange() {
	if err := r.channel.channel.ExchangeDeclare(
		r.channel.exchangeName,
		r.kind,
		r.durable,
		r.autoDelete,
		r.internal,
		r.noWait,
		nil,
	); err != nil {
		r.logger.Fatalf("failed to create exchange: '%s'", r.channel.exchangeName)
	}
}

func (r *RmqClient) BindQueue() {
	if err := r.channel.channel.QueueBind(
		r.channel.queueName,
		r.channel.routingKey,
		r.channel.exchangeName,
		r.noWait,
		nil,
	); err != nil {
		r.logger.Fatalf("failed to bind queue and exchange: '%s'", r.channel.queueName)
	}
}

func (r *RmqClient) CloseConnection() {
	err := r.connection.Close()
	if err != nil {
		r.logger.Errorf("failed to close connection")
	}
}

func (r *RmqClient) CloseChannel() {
	err := r.channel.channel.Close()
	if err != nil {
		r.logger.Errorf("failed to close channel")
	}
}
