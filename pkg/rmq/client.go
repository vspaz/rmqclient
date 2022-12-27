package rmq

import (
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"time"
)

type Client struct {
	connectionUrl string
	heartBeat     time.Duration

	connection *amqp.Connection
	logger     *logrus.Logger
}

func NewClient(connectionUrl string, logger *logrus.Logger) *Client {
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

func (c *Client) CloseConnection() {
	err := c.connection.Close()
	if err != nil {
		c.logger.Errorf("failed to close connection")
	}
}
