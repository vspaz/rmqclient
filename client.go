package rmq

import (
	"bytes"
	"fmt"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"net/http"
	"os"
	"time"
)

type RQClient struct {
	host     string
	user     string
	password string
	port     int

	kind       string
	durable    bool
	autoDelete bool
	internal   bool
	noWait     bool
	exclusive  bool

	exchangeName string
	queueName    string
	heartBeat    time.Duration

	logger     *zap.SugaredLogger
	httpclient http.Client
}

func New(server string, user string, password string, port int, logger *zap.SugaredLogger, httpclient http.Client) *RQClient {
	return &RQClient{
		host:     server,
		user:     user,
		password: password,
		port:     port,

		kind:       "direct",
		durable:    true,
		autoDelete: false,
		internal:   false,
		noWait:     false,
		exclusive:  true,

		exchangeName: "task.results",
		queueName:    "task.results",
		heartBeat:    60,

		logger:     logger,
		httpclient: httpclient,
	}
}

func (rq *RQClient) connect() *amqp.Connection {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d", rq.user, rq.password, rq.host, rq.port)
	connection, err := amqp.DialConfig(url, amqp.Config{Heartbeat: time.Second * rq.heartBeat})
	if err != nil {
		rq.logger.Error("failed to establish connection at ", url)
		os.Exit(-1)
	}
	return connection
}

func (rq *RQClient) createChannel(connection *amqp.Connection) *amqp.Channel {
	channel, err := connection.Channel()
	if err != nil {
		rq.logger.Error("failed to create channel")
		os.Exit(-1)
	}
	return channel
}

func (rq *RQClient) declareExchange(chanel *amqp.Channel) {
	if err := chanel.ExchangeDeclare(
		rq.exchangeName,
		rq.kind,
		rq.durable,
		rq.autoDelete,
		rq.internal,
		rq.noWait,
		nil,
	); err != nil {
		rq.logger.Error("failed to create exchange: ", rq.exchangeName)
		os.Exit(-1)
	}
}

func (rq *RQClient) declareQueue(channel *amqp.Channel) {
	if _, err := channel.QueueDeclare(
		rq.queueName,
		rq.durable,
		rq.autoDelete,
		rq.exclusive,
		rq.noWait,
		nil,
	); err != nil {
		rq.logger.Error("failed to declare queue: ", rq.queueName)
		os.Exit(-1)
	}
}

func (rq *RQClient) bindQueue(channel *amqp.Channel) {
	if err := channel.QueueBind(
		rq.queueName,
		"",
		rq.exchangeName,
		rq.noWait,
		nil,
	); err != nil {
		rq.logger.Error("failed to bind queue and exchange: ", rq.queueName)
		os.Exit(-1)
	}
}

func (rq *RQClient) consumeTasks(channel *amqp.Channel) <-chan amqp.Delivery {
	resultChannel, err := channel.Consume(
		rq.queueName,
		"notifier",
		true,
		rq.exclusive,
		false,
		rq.noWait,
		nil,
	)
	if err != nil {
		rq.logger.Error("failed to consume tasks")
		os.Exit(-1)
	}
	return resultChannel
}

func (rq *RQClient) notify(url string, contentType string, body []byte) {
	resp, err := rq.httpclient.Post(url, contentType, bytes.NewBuffer(body))
	if err != nil || resp.StatusCode != 202 {
		rq.logger.Error("error to notify; received status code: ", resp.StatusCode)
	}
	defer resp.Body.Close()
}

func (rq *RQClient) processTasks(Channel <-chan amqp.Delivery) {
	for result := range Channel {
		fmt.Println(result.Body)
	}
}
