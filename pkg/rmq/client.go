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

	queueName string
	heartBeat time.Duration

	channel    *amqp.Channel
	connection *amqp.Connection
	logger     *logrus.Logger
}

func New(connectionUrl string, queueName string, logger *logrus.Logger) *RmqClient {
	return &RmqClient{
		connectionUrl: connectionUrl,
		kind:          "direct",
		durable:       true,
		autoDelete:    false,
		internal:      false,
		noWait:        false,
		exclusive:     false,

		queueName: queueName,
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

func (r *RmqClient) CreateChannel() {
	r.logger.Info("trying to create a channel")
	channel, err := r.connection.Channel()
	if err != nil {
		r.logger.Fatalf("failed to create channel")
	}
	r.logger.Info("channel created: OK")
	r.channel = channel
}

func (r *RmqClient) DeclareExchange(exchangeName string) {
	if err := r.channel.ExchangeDeclare(
		exchangeName,
		r.kind,
		r.durable,
		r.autoDelete,
		r.internal,
		r.noWait,
		nil,
	); err != nil {
		r.logger.Fatalf("failed to create exchange: '%s'", exchangeName)
	}
}

func (r *RmqClient) BindQueue(exchangeName, routingKey string) {
	if err := r.channel.QueueBind(
		r.queueName,
		routingKey,
		exchangeName,
		r.noWait,
		nil,
	); err != nil {
		r.logger.Fatalf("failed to bind queue and exchange: '%s'", r.queueName)
	}
}

func (r *RmqClient) CloseConnection() {
	err := r.connection.Close()
	if err != nil {
		r.logger.Errorf("failed to close connection")
	}
}

func (r *RmqClient) CloseChannel() {
	err := r.channel.Close()
	if err != nil {
		r.logger.Errorf("failed to close channel")
	}
}
