package rmq

import (
	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"time"
)

type Connection struct {
	connectionUrl string
	heartBeat     time.Duration
	connection    *amqp091.Connection
	logger        *logrus.Logger
}

func NewConnection(connectionUrl string, logger *logrus.Logger) *Connection {
	return &Connection{
		connectionUrl: connectionUrl,
		heartBeat:     60 * time.Second,
		logger:        logger,
	}
}

func (c *Connection) Create() {
	c.logger.Debugf("connecting to rabbitmq '%s'", c.connectionUrl)
	connection, err := amqp091.DialConfig(c.connectionUrl, amqp091.Config{Heartbeat: time.Second * c.heartBeat})
	if err != nil {
		c.logger.Fatalf("failed to establish connection at '%s'", c.connectionUrl)
	}
	c.logger.Info("connection to rabbitmq established at: OK")
	c.connection = connection
}

func (c *Connection) Close() {
	err := c.connection.Close()
	if err != nil {
		c.logger.Errorf("failed to close connection")
	}
}
